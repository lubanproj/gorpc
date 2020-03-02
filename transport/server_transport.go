package transport

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/protocol"
	"github.com/lubanproj/gorpc/stream"
	"github.com/lubanproj/gorpc/utils"
	"io"
	"net"
	"time"
)

type serverTransport struct {
	opts *ServerTransportOptions
}

var serverTransportMap = make(map[string]ServerTransport)

func init() {
	serverTransportMap["default"] = DefaultServerTransport
}

func GetServerTransport(transport string) ServerTransport {

	if v, ok := serverTransportMap[transport]; ok {
		return v
	}

	return DefaultServerTransport
}

var DefaultServerTransport = NewServerTransport()

var NewServerTransport = func() ServerTransport {
	return &serverTransport{
		opts : &ServerTransportOptions{},
	}
}

func (s *serverTransport) ListenAndServe(ctx context.Context, opts ...ServerTransportOption) error {

	for _, o := range opts {
		o(s.opts)
	}

	switch s.opts.Network {
		case "tcp","tcp4","tcp6":
			return s.ListenAndServeTcp(ctx, opts ...)
		case "udp","udp4", "udp6":
			return s.ListenAndServeUdp(ctx, opts ...)
		default:
			return codes.NetworkNotSupportedError
	}
}

func (s *serverTransport) ListenAndServeTcp(ctx context.Context, opts ...ServerTransportOption) error {

	lis, err := net.Listen(s.opts.Network, s.opts.Address)
	if err != nil {
		return err
	}

	var tempDelay time.Duration
	for {

		tl, ok := lis.(*net.TCPListener);
		if !ok {
			return codes.NetworkNotSupportedError
		}

		conn , err := tl.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		if err = conn.SetKeepAlive(true); err != nil {
			return err
		}

		if s.opts.KeepAlivePeriod != 0 {
			conn.SetKeepAlivePeriod(s.opts.KeepAlivePeriod)
		}

		go func() {

			// build stream
			ctx, _ := stream.NewServerStream(ctx)

			if err := s.handleConn(ctx, wrapConn(conn)); err != nil {
				log.Error("gorpc handle tcp conn error, %v", err)
			}

		}()

	}

}


func (s *serverTransport) handleConn(ctx context.Context, conn *connWrapper) error {

	// close the connection before return
	defer conn.Close()

	for {
		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err();
		default:
		}

		frame , err := s.read(ctx, conn)
		if err == io.EOF {
			// read compeleted
			return nil
		}

		if err != nil {
			return err
		}

		// parse protocol header
		request := &protocol.Request{}
		if err = proto.Unmarshal(frame[codec.FrameHeadLen:], request); err != nil {
			return err
		}

		// build serverStream
		_, err = s.getServerStream(ctx, request)
		if err != nil {
			return err
		}

		rsp , err := s.handle(ctx, request.Payload)
		if err != nil {
			return err
		}

		if err = s.write(ctx, conn,rsp); err != nil {
			return err
		}
	}

}

func (s *serverTransport) read(ctx context.Context, conn *connWrapper) ([]byte, error) {

	frame, err := conn.framer.ReadFrame(conn)

	if err != nil {
		return nil, err
	}

	return frame, nil
}


func (s *serverTransport) handle(ctx context.Context, payload []byte) ([]byte, error) {

	rsp , err := s.opts.Handler.Handle(ctx, payload)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *serverTransport) write(ctx context.Context, conn net.Conn, rsp []byte) error {
	_, err := conn.Write(rsp)

	return err
}


type connWrapper struct {
	net.Conn
	framer Framer
}

func wrapConn(rawConn net.Conn) *connWrapper {
	return &connWrapper{
		Conn : rawConn,
		framer: NewFramer(),
	}
}


func (s *serverTransport) getServerStream(ctx context.Context, request *protocol.Request) (*stream.ServerStream, error) {
	serverStream := stream.GetServerStream(ctx)

	_, method , err := utils.ParseServicePath(string(request.ServicePath))
	if err != nil {
		return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	}

	serverStream.WithMethod(method)

	return serverStream, nil
}


