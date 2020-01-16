package transport

import (
	"context"
	"time"
)

type ServerTransportOptions struct{
	Address string // 地址，格式例如 ip://127.0.0.1：8080
	Network string  // 网络类型
	Timeout time.Duration  // 传输层请求超时时间，默认为 2 min
	Handler Handler
	Serialization string   // 序列化方式
}

type Handler interface {
	Handle (context.Context, []byte) ([]byte, error)
}


type ServerTransportOption func(*ServerTransportOptions)

func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Address = address
	}
}

func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Network = network
	}
}

func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Timeout = timeout
	}
}

func WithHandler(handler Handler) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Handler = handler
	}
}

func WithSerialization(serialization string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Serialization = serialization
	}
}