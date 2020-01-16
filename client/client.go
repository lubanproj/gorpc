package client

import (
	"context"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/pool/connpool"
	"github.com/lubanproj/gorpc/stream"
	"github.com/lubanproj/gorpc/transport"
	"strings"
)

// Client 定义了客户端通用接口
type Client interface {
	Invoke(ctx context.Context, req , rsp interface{}, path string, opts ...Option) error
}

// 全局使用一个 client
var DefaultClient = New()

var New = func() Client {
	return &defaultClient{
		opts : &Options{
			protocol : "proto",
		},
	}
}

type defaultClient struct {
	opts *Options
}

func (c *defaultClient) Invoke(ctx context.Context, req , rsp interface{}, path string, opts ...Option) error {

	for _, opt := range opts {
		opt(c.opts)
	}

	// 设置服务名、方法名
	newCtx, clientStream := stream.NewClientStream(ctx)

	index := strings.LastIndex(path, "/")
	if index == 0 {
		return codes.NewFrameworkError(codes.ClientDialErrorCode, "invalid path")
	}
	c.opts.serviceName = path[1:index]
	c.opts.method = path[index+1:]

	// 这里先保留看看，需不需要去掉
	clientStream.WithServiceName(path[1:index])
	clientStream.WithMethod(path[index+1:])

	// 先执行拦截器
	return interceptor.Intercept(newCtx, req, rsp, c.opts.interceptors, c.invoke)
}

func (c *defaultClient) invoke(ctx context.Context, req, rsp interface{}) error {

	serialization := codec.GetSerialization(c.opts.protocol)
	reqbuf, err := serialization.Marshal(req)
	if err != nil {
		return codes.ClientMsgError
	}

	clientCodec := codec.GetCodec(c.opts.protocol)
	reqbody, err := clientCodec.Encode(reqbuf)
	if err != nil {
		return err
	}

	clientTransport := c.NewClientTransport()
	clientTransportOpts := []transport.ClientTransportOption {
		transport.WithClientTarget(c.opts.target),
		transport.WithClientNetwork(c.opts.network),
		transport.WithClientPool(connpool.GetPool("default")),
	}
	rspbuf, err := clientTransport.Send(ctx, reqbody, clientTransportOpts ...)
	if err != nil {
		return err
	}

	rspbody, err := clientCodec.Decode(rspbuf)
	if err != nil {
		return err
	}

	return serialization.Unmarshal(rspbody, rsp)

}

func (c *defaultClient) NewClientTransport() transport.ClientTransport {
	return transport.GetClientTransport(c.opts.protocol)
}







