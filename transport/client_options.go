package transport

import (
	"github.com/lubanproj/gorpc/pool/connpool"
	"github.com/lubanproj/gorpc/selector"
)

type ClientTransportOptions struct {
	Target string
	ServiceName string
	Network string
	Pool connpool.Pool
	Selector selector.Selector
}

type ClientTransportOption func(*ClientTransportOptions)

func WithServiceName(serviceName string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.ServiceName = serviceName
	}
}

func WithClientTarget(target string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Target = target
	}
}

func WithClientNetwork(network string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Network = network
	}
}

func WithClientPool(pool connpool.Pool) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Pool = pool
	}
}

func WithSelector(selector selector.Selector) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Selector = selector
	}
}
