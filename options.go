package gorpc

import (
	"time"
)

type ServerOptions struct {
	address string  // 监听地址，格式 ip://127.0.0.1:8080 , dns://www.google.com
	network string  // 网络类型，如 tcp、udp
	protocol string  // 协议类型，如 proto/json 等
	timeout time.Duration       // 超时时间
	serialization string 	// 序列化方式，默认是 proto

	consulAddr string       // consul server 地址，当服务发现方式为 consul 时需要填写
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption{
	return func(o *ServerOptions) {
		o.address = address
	}
}

func WithNetwork(network string) ServerOption {
	return func(o *ServerOptions) {
		o.network = network
	}
}

func WithProtocol(protocol string) ServerOption {
	return func(o *ServerOptions) {
		o.protocol = protocol
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}

func WithSerialization(serialization string) ServerOption {
	return func(o *ServerOptions) {
		o.serialization = serialization
	}
}

func WithConsulAddr(consulAddr string) ServerOption {
	return func(o *ServerOptions) {
		o.consulAddr = consulAddr
	}
}

