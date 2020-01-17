package transport

import (
	"context"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/protocol"
	"github.com/lubanproj/gorpc/stream"
	"github.com/lubanproj/gorpc/utils"
	"net"
)

const GORPCHeaderLength = 5

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
		return codes.NewFrameworkError(codes.ServerNetworkErrorCode, err.Error())
	}

	for {
		conn , err := lis.Accept()

		if err != nil {
			return codes.NewFrameworkError(codes.ServerNetworkErrorCode, err.Error())
		}

		go func() {
			// 构造 stream
			newCtx, _ := stream.NewServerStream(ctx)

			if err := s.handleConn(newCtx, conn); err != nil {
				log.Error("gorpc handle conn error, %v", err)
			}
		}()

	}

}

func (s *serverTransport) ListenAndServeUdp(ctx context.Context, opts ...ServerTransportOption) error {

	return nil
}

func (s *serverTransport) handleConn(ctx context.Context, rawConn net.Conn) error {

	// rawConn.SetDeadline(time.Now().Add(s.opts.Timeout))
	// tcpConn := newTcpConn(rawConn)
	frame , err := s.read(ctx,rawConn)
	if err != nil {
		return err
	}

	// 解析协议头
	ser := codec.GetSerialization(s.opts.Serialization)
	request := &protocol.Request{}

	if err = ser.Unmarshal(frame[codec.FrameHeadLen:], request); err != nil {
		return err
	}

	// 构造 serverStream
	_, err = s.buildServerStream(ctx, request)
	if err != nil {
		return err
	}

	rsp , err := s.handle(ctx, frame)
	if err != nil {
		return err
	}

	err = s.write(ctx, rawConn,rsp)
	return err
}

func (s *serverTransport) read(ctx context.Context, conn net.Conn) ([]byte, error) {

	frame, err := codec.ReadFrame(conn)

	if err != nil {
		return nil, err
	}

	return frame, nil
}


func (s *serverTransport) handle(ctx context.Context, frame []byte) ([]byte, error) {

	rsp , err := s.opts.Handler.Handle(ctx, frame)
	if err != nil {
		return nil, codes.NewFrameworkError(codes.ServerNoResponseErrorCode, err.Error())
	}

	return rsp, nil
}

func (s *serverTransport) write(ctx context.Context, conn net.Conn, rsp []byte) error {
	_, err := conn.Write(rsp)

	return err
}


type tcpConn struct {
	conn net.Conn

}

func newTcpConn(rawConn net.Conn) *tcpConn {
	return &tcpConn{
		conn : rawConn,
	}
}


func (s *serverTransport) buildServerStream(ctx context.Context, request *protocol.Request) (*stream.ServerStream, error) {
	serverStream := stream.GetServerStream(ctx)

	_, method , err := utils.ParseServicePath(string(request.ServicePath))
	if err != nil {
		return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	}

	serverStream.WithMethod(method)

	return serverStream, nil
}


