package transport

import (
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/pool/connpool"
)

type ClientTransportOptions struct {
	Target string
	Network string
	Codec codec.Codec
	Serialization codec.Serialization
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

func WithClientCodec(codec codec.Codec) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Codec = codec
	}
}

func WithClientSerialization(serialization codec.Serialization) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Serialization = serialization
	}
}

func WithClientPool(pool connpool.Pool) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Pool = pool
	}
}
