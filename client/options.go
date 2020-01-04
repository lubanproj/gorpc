package client

import (
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/transport"
	"time"
)

// Options 定义了客户端调用参数
type Options struct {
	target string 	// 调用地址，格式为 ip:port 127.0.0.1:8000
	timeout time.Duration 	// 超时时间
	network string  // 网络类型 tcp/udp
	protocol   string  // 协议类型 proto/json
	serializedType string // 序列化类型
	transportOpts transport.ClientTransportOptions
	interceptors []interceptor.ClientInterceptor
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

func WithNetwork(network string) Option {
	return func(o *Options) {
		o.network = network
	}
}

func WithProtocol(protocol string) Option {
	return func(o *Options) {
		o.protocol = protocol
	}
}

func WithSerializedType(serializedType string) Option {
	return func(o *Options) {
		o.serializedType = serializedType
	}
}