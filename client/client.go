package client

import (
	"context"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/interceptor"
)

// 全局使用一个 client
var DefaultClient = New()

var New = func() Client {
	return &defaultClient{}
}

// Client 定义了客户端通用接口
type Client interface {
	Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error
}


type defaultClient struct {
	options *Options
}

func (c *defaultClient) Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error {
	for _, opt := range opts {
		opt(c.options)
	}

	// 先执行拦截器
	return interceptor.Intercept(ctx, req, rsp, c.options.interceptors, c.invoke)
}

func (c *defaultClient) invoke(ctx context.Context, req interface{}, rsp interface{}) error {

	reqBytes, ok := req.([]byte)
	if !ok {
		return codes.ClientMsgError
	}

	if _, err := c.options.Transport.Send(ctx, reqBytes); err != nil {
		return codes.ClientNetworkError
	}

	return nil
}



