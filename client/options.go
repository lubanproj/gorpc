package client

import (
	"time"

	"github.com/lubanproj/gorpc/auth"
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/transport"
)

// Options defines the client call parameters
type Options struct {
	serviceName string // service name
	method string // method name
	target string 	// format e.g.:  ip:port 127.0.0.1:8000
	timeout time.Duration  // timeout
	network string  // network type, e.g.:  tcp、udp
	protocol   string  // protocol type , e.g. : proto、json
	serializationType string // seralization type , e.g. : proto、msgpack
	transportOpts transport.ClientTransportOptions
	interceptors []interceptor.ClientInterceptor
	selectorName string      // service discovery name, e.g. : consul、zookeeper、etcd
	perRPCAuth []auth.PerRPCAuth  // authentication information required for each RPC call
	transportAuth auth.TransportAuth
}

type Option func(*Options)

func WithServiceName(serviceName string) Option {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

func WithMethod(method string) Option {
	return func(o *Options) {
		o.method = method
	}
}

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

func WithSerializationType(serializationType string) Option {
	return func(o *Options) {
		o.serializationType = serializationType
	}
}

func WithSelectorName(selectorName string) Option {
	return func(o *Options) {
		o.selectorName = selectorName
	}
}

func WithInterceptor(interceptors ...interceptor.ClientInterceptor) Option {
	return func(o *Options) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithPerRPCAuth(rpcAuth auth.PerRPCAuth) Option {
	return func(o *Options) {
		o.perRPCAuth = append(o.perRPCAuth,rpcAuth)
	}
}

func WithTransportAuth(transportAuth auth.TransportAuth) Option {
	return func(o *Options) {
		o.transportAuth = transportAuth
	}
}

