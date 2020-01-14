package transport

import (
	"github.com/lubanproj/gorpc/pool/connpool"
)

type ClientTransportOptions struct {
	Target string
	Network string
	Pool connpool.Pool
}

type ClientTransportOption func(*ClientTransportOptions)

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
