package transport

import (
	"context"
	"time"
)

type ServerTransportOptions struct{
	Address string // address，e.g: ip://127.0.0.1：8080
	Network string  // network type
	Timeout time.Duration  // Transport layer request timeout ，default: 2 min
	Handler Handler		   // handler
	Serialization string   // serialization type, e.g : proto、json、msgpack
	KeepAlivePeriod time.Duration // keepalive period
}

type Handler interface {
	Handle (context.Context, []byte) ([]byte, error)
}


type ServerTransportOption func(*ServerTransportOptions)

func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Address = address
	}
}

func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Network = network
	}
}

func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Timeout = timeout
	}
}

func WithHandler(handler Handler) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Handler = handler
	}
}

func WithSerialization(serialization string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Serialization = serialization
	}
}

func WithKeepAlivePeriod(keepAlivePeriod time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.KeepAlivePeriod = keepAlivePeriod
	}
}