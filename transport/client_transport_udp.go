package transport

import (
	"context"
	"github.com/lubanproj/gorpc/codes"
	"net"
)

func (c *clientTransport) SendUdpReq(ctx context.Context, req []byte) ([]byte, error) {
	// service discovery
	addr, err := c.opts.Selector.Select(c.opts.ServiceName)
	if err != nil {
		return nil, err
	}

	// defaultSelector returns "", use the target as address
	if addr == "" {
		addr = c.opts.Target
	}

	udpAddr, err := net.ResolveUDPAddr(c.opts.Network, addr)
	if err != nil {
		return nil, codes.NewFrameworkError(codes.ClientMsgErrorCode, "addr invalid ...")
	}

	conn, err := net.DialUDP(c.opts.Network, nil, udpAddr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if n, err := conn.Write(req); n != len(req) || err != nil {
		return nil, err
	}

	recvBuf := make([]byte, 65536)
	n, err := conn.Read(recvBuf);
	if err != nil {
		return nil, err
	}

	rsp := recvBuf[:n]

	return rsp, nil
}