package transport

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/lubanproj/gorpc/codec"
	"io"
	"net"
	"time"

	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/protocol"
	"github.com/lubanproj/gorpc/stream"
	"github.com/lubanproj/gorpc/utils"
)

type serverTransport struct {
	opts *ServerTransportOptions
}

var serverTransportMap = make(map[string]ServerTransport)

func init() {
	serverTransportMap["default"] = DefaultServerTransport
}

// RegisterServerTransport supports business custom registered ServerTransport
func RegisterServerTransport(name string, serverTransport ServerTransport) {
	if serverTransportMap == nil {
		serverTransportMap = make(map[string]ServerTransport)
	}
	serverTransportMap[name] = serverTransport
}

// Get the ServerTransport
func GetServerTransport(transport string) ServerTransport {

	if v, ok := serverTransportMap[transport]; ok {
		return v
	}

	return DefaultServerTransport
}

// The default server transport
var DefaultServerTransport = NewServerTransport()

// Use the singleton pattern to create a server transport
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

	go func() {
		if err = s.serve(ctx ,lis); err != nil {
			log.Errorf("transport serve error, %v", err)
		}
	}()

	return nil
}

func (s *serverTransport) serve(ctx context.Context,lis net.Listener) error {

	var tempDelay time.Duration

	tl, ok := lis.(*net.TCPListener)
	if !ok {
		return codes.NetworkNotSupportedError
	}

	for {

		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err();
		default:
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
				log.Errorf("gorpc handle tcp conn error, %v", err)
			}

		}()

	}

	return nil
}


func (s *serverTransport) handleConn(ctx context.Context, conn *connWrapper) error {

	// close the connection before return
	// the connection closes only if a network read or write fails
	defer conn.Close()

	for {
		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err()
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

		rsp , err := s.handle(ctx, frame)
		if err != nil {
			log.Errorf("s.handle err is not nil, %v", err)
		}

		if err = s.write(ctx, conn, rsp); err != nil {
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


func (s *serverTransport) handle(ctx context.Context, frame []byte) ([]byte, error) {

	// parse reqbuf into req interface {}
	serverCodec := codec.GetCodec(s.opts.Protocol)

	reqbuf, err := serverCodec.Decode(frame)
	if err != nil {
		log.Errorf("server Decode error: %v", err)
		return nil, err
	}

	rspbuf, err := s.opts.Handler.Handle(ctx, reqbuf)
	if err != nil {
		log.Errorf("server Handle error: %v", err)
	}

	response := addRspHeader(rspbuf, err)

	rspPb, err := proto.Marshal(response)
	if err != nil {
		log.Errorf("proto Marshal error: %v", err)
		return nil, err
	}

	rspbody, err := serverCodec.Encode(rspPb)
	if err != nil {
		log.Errorf("server Encode error, response: %v, err: %v", response, err)
		return nil, err
	}

	log.Debugf("rspbody : %v, err : %v", rspbody, err)

	return rspbody, nil
}

func addRspHeader(payload []byte, err error) *protocol.Response {
	response := &protocol.Response{
		Payload: payload,
		RetCode: codes.OK,
		RetMsg: "success",
	}

	if err != nil {
		if e, ok := err.(*codes.Error); ok {
			response.RetCode = e.Code
			response.RetMsg = e.Message
		} else {
			response.RetCode = codes.ServerInternalErrorCode
			response.RetMsg = codes.ServerInternalError.Message
		}
	}

	return response
}

func (s *serverTransport) write(ctx context.Context, conn net.Conn, rsp []byte) error {
	if _, err := conn.Write(rsp); err != nil {
		log.Errorf("conn Write err: %v", err)
	}

	return nil
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


