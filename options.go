package gorpc

import "time"

type ServerOptions struct {
	address string  // 监听地址，格式 ip://127.0.0.1:8080 , dns://www.google.com
	network string  // 网络类型，如 tcp、udp
	protocol string  // 协议类型，如 proto/json 等
	method  string   // 请求的方法名
	timeout time.Duration       // 超时时间
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

func WithMethod(method string) ServerOption {
	return func(o *ServerOptions) {
		o.method = method
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}



