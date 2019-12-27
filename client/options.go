package client

import (
	"github.com/diubrother/gorpc/interceptor"
	"github.com/diubrother/gorpc/transport"
	"time"
)

// Options 定义了客户端调用参数
type Options struct {
	// 调用地址
	target string
	// 超时时间
	timeout time.Duration

	Transport transport.ClientTransport

	interceptors []interceptor.Interceptor
}

type Option func(*Options)

func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}