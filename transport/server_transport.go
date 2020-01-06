package transport

import (
	"context"
	"encoding/binary"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/log"
	"golang.org/x/net/http2"
	"io"
	"net"
	"time"
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

		go s.handleConn(ctx , conn)


	}
	return nil
}

func (s *serverTransport) ListenAndServeUdp(ctx context.Context, opts ...ServerTransportOption) error {

	return nil
}

func (s *serverTransport) handleConn(ctx context.Context, rawConn net.Conn) error {

	rawConn.SetDeadline(time.Now().Add(s.opts.Timeout))
	tcpConn := newTcpConn(rawConn)
	req , err := s.read(ctx,tcpConn)

	if err != nil {
		return err
	}

	rsp , err := s.handle(ctx,tcpConn, req)
	if err != nil {
		return err
	}

	err = s.write(ctx,tcpConn,rsp)
	return err
}

func (s *serverTransport) read(ctx context.Context, conn *tcpConn) ([]byte, error) {
	// 先读出 http包头
	http2.ReadFrameHeader(conn.conn)

	// 再读出协议包头
	header := make([]byte, GORPCHeaderLength)
	io.ReadFull(conn.conn, header)

	compressingType := header[0]

	if compressingType == 1 {
		// TODO 压缩模式，需要解压缩
	}

	length := binary.BigEndian.Uint32(header[1:])
	msg := make([]byte, length)
	_, err := io.ReadFull(conn.conn, msg)

	if err != nil {
		log.Error("read data from conn error, %v", err)
		return nil, codes.ServerDecodeError
	}
	return msg, nil
}

func (s *serverTransport) handle(ctx context.Context, conn *tcpConn, req []byte) ([]byte, error) {

	rsp , err := s.opts.Handler(ctx, req)
	if err != nil {
		return nil, err
	}

	rspbuf, err := s.opts.Serialization.Marshal(rsp)
	if err != nil {
		return nil, err
	}

	rspbody, err := s.opts.Codec.Encode(rspbuf)
	if err != nil {
		return nil, err
	}

	return rspbody, nil
}

func (s *serverTransport) write(ctx context.Context, conn *tcpConn, rsp []byte) error {
	_, err := conn.conn.Write(rsp)

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



