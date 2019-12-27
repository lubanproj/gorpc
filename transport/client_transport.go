package transport

import (
	"context"
)

type clientTransport struct {
	opts *ClientTransportOptions
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
			return nil, codes.NewFrameworkError(codes.ClientNetworkErrorCode,"conn send error")
		}
		sendNum += num

		if err = isDone(ctx); err != nil {
			return nil, err
		}
	}

	// 解包

	return nil, nil
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