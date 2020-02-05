package gorpc

import (
	"github.com/lubanproj/gorpc/interceptor"
	"time"
)

type ServerOptions struct {
	address string  // 监听地址，格式 ip://127.0.0.1:8080 , dns://www.google.com
	network string  // 网络类型，如 tcp、udp
	protocol string  // 协议类型，如 proto/json 等
	timeout time.Duration       // 超时时间
	serializationType string 	// 序列化方式，默认是 proto

	selectorSvrAddr string       // 服务发现 server 地址，当使用第三方服务发现方式时需要填写
	tracingSvrAddr  string 		 // tracing 类插件 server 地址，当使用第三方 tracing 类插件时需要填写
	tracingSpanName string       // tracing 类插件 span name, 当使用第三方 tracing 类插件时需要填写
	pluginNames []string         // 插件名字
	interceptors []interceptor.ServerInterceptor
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

func WithSerializationType(serializationType string) ServerOption {
	return func(o *ServerOptions) {
		o.serializationType = serializationType
	}
}

func WithSelectorSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.selectorSvrAddr = addr
	}
}

func WithPlugin(pluginName ... string) ServerOption {
	return func(o *ServerOptions) {
		o.pluginNames = append(o.pluginNames, pluginName ...)
	}
}

func WithInterceptor(interceptors ...interceptor.ServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithTracingSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSvrAddr = addr
	}
}

func WithTracingSpanName(name string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSpanName = name
	}
}