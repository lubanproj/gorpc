package transport

import (
	"context"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/pool/connpool"
)

type clientTransport struct {
	opts *ClientTransportOptions
}

var DefaultClientTransport = New()

var New = func() ClientTransport {
	return &clientTransport{
		opts : &ClientTransportOptions{
			NetworkType: "tcp",
			Codec : codec.DefaultCodec,
			Serialization: codec.DefaultSerialization,
		},
	}
}

func (c *clientTransport) Send(ctx context.Context, req []byte, opts ...ClientTransportOption) ([]byte, error) {
	if c.opts.NetworkType == "tcp" {
		return c.SendTcpReq(ctx, req)
	}

	if c.opts.NetworkType == "udp" {
		return c.SendUdpReq(ctx, req)
	}

	return nil, codes.NetworkNotSupportedError
}

func (c *clientTransport) SendTcpReq(ctx context.Context, req []byte) ([]byte, error) {

	// 从连接池里面获取一个连接
	conn, err := connpool.DefaultPool.Get(ctx, "tcp", c.opts.Target)
	if err != nil {
		return nil, codes.ConnectionError
	}
	defer conn.Close()

	sendNum := 0
	num := 0
	for sendNum < len(req) {
		num , err = conn.Write(req)
		if err != nil {
			return nil, codes.NewFrameworkError(codes.ClientNetworkErrorCode,err.Error())
		}
		sendNum += num

		if err = isDone(ctx); err != nil {
			return nil, err
		}
	}

	// 解析帧
	rspbuf, err := codec.ReadFrame(conn)
	if err != nil {
		return nil, codes.NewFrameworkError(codes.ClientNetworkErrorCode, err.Error())
	}

	rspbody, err := c.opts.Codec.Decode(rspbuf)

	return rspbody, err
}

func (c *clientTransport) SendUdpReq(ctx context.Context, req []byte) ([]byte, error) {

	return nil, nil
}


func isDone(ctx context.Context) error {
	select {
	case <- ctx.Done() :
		if ctx.Err() == context.Canceled {
			return codes.ClientContextCanceledError
		}
		if ctx.Err() == context.DeadlineExceeded {
			return codes.ClientTimeoutError
		}
	default:
	}

	return nil
}