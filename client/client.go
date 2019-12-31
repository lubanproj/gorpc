package client

import (
	"context"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/interceptor"
)

// 全局使用一个 client
var DefaultClient = New()

var New = func() Client {
	return &defaultClient{
		opts : &Options{},
	}
}

// Client 定义了客户端通用接口
type Client interface {
	Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error
}


type defaultClient struct {
	opts *Options
}

func (c *defaultClient) Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error {

	for _, opt := range opts {
		opt(c.opts)
	}

	// 先执行拦截器
	return interceptor.Intercept(ctx, req, rsp, c.opts.interceptors, c.invoke)
}

func (c *defaultClient) invoke(ctx context.Context, req,rsp interface{}) error {

	reqBytes, ok := req.([]byte)
	if !ok {
		return codes.ClientMsgError
	}

	rsp, err := c.opts.transport.Send(ctx, reqBytes)

	if err != nil {
		return codes.ClientNetworkError
	}

	return nil
}



