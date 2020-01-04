package gorpc

type ServerOptions struct {
	target string  // 监听地址，格式 ip://127.0.0.1:8080 , dns://www.google.com
	network string  // 网络类型，如 tcp、udp
	protocol string  // 协议类型，如 proto/json 等
}

type ServerOption func(*ServerOptions)

func WithTarget(target string) ServerOption{
	return func(o *ServerOptions) {
		o.target = target
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

type ServiceOptions struct {
	target string  // 监听地址，格式 ip://127.0.0.1:8080 , dns://www.google.com
	network string  // 网络类型，如 tcp、udp
	protocol string  // 协议类型，如 proto/json 等
}

