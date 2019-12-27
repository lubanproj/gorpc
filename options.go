package gorpc

import (
	"github.com/luban_proj/gorpc/transport"
)

type ServiceOptions struct {
	target string  // 监听地址，格式 ip://127.0.0.1:8080 , dns://www.google.com
	network string  // 网络类型，如 tcp、udp
	transport transport.ServerTransport  // server 端传输层
	transportOptions []transport.ServerTransportOption  // server 端传输层参数选项
}

type ServiceOption func(*ServiceOptions)

func WithTarget(target string) ServiceOption{
	return func(o *ServiceOptions) {
		o.target = target
	}
}

