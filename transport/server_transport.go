package transport

import (
	"context"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/log"
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
			log.Error("sssssssssssssssss")
			if err := s.handleConn(ctx, conn); err != nil {
				log.Error("gorpc handle conn error, %v", err)
			}
		}()

	}
	return nil
}

func (s *serverTransport) ListenAndServeUdp(ctx context.Context, opts ...ServerTransportOption) error {

	return nil
}

func (s *serverTransport) handleConn(ctx context.Context, rawConn net.Conn) error {

	// rawConn.SetDeadline(time.Now().Add(s.opts.Timeout))
	// tcpConn := newTcpConn(rawConn)
	req , err := s.read(ctx,rawConn)

	if err != nil {
		return err
	}

	rsp , err := s.handle(ctx, req)
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


func (s *serverTransport) handle(ctx context.Context, req []byte) ([]byte, error) {

	rsp , err := s.opts.Handler.Handle(ctx, req)
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



