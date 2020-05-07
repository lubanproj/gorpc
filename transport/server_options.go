package transport

import (
	"context"
	"time"
)

// ServerTransportOptions includes all ServerTransport parameter options
type ServerTransportOptions struct{
	Address string // address，e.g: ip://127.0.0.1：8080
	Network string  // network type
	Protocol string  // protocol type, e.g. : proto、json
	Timeout time.Duration  // transport layer request timeout ，default: 2 min
	Handler Handler		   // handler
	Serialization string   // serialization type, e.g : proto、json、msgpack
	KeepAlivePeriod time.Duration // keepalive period
}

// Handler defines a common interface for handling packets
type Handler interface {
	Handle (context.Context, []byte) ([]byte, error)
}

// Use the Options mode to wrap the ServerTransportOptions
type ServerTransportOption func(*ServerTransportOptions)

// WithServerAddress returns a ServerTransportOption which sets the value for address
func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Address = address
	}
}

// WithServerNetwork returns a ServerTransportOption which sets the value for network
func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Network = network
	}
}

// WithProtocol returns a ServerTransportOption which sets the value for protocol
func WithProtocol(protocol string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Protocol = protocol
	}
}

// WithServerTimeout returns a ServerTransportOption which sets the value for timeout
func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Timeout = timeout
	}
}

// WithHandler returns a ServerTransportOption which sets the value for handler
func WithHandler(handler Handler) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Handler = handler
	}
}

// WithSerialization returns a ServerTransportOption which sets the value for serialization
func WithSerialization(serialization string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Serialization = serialization
	}
}

// WithKeepAlivePeriod returns a ServerTransportOption which sets the value for keepAlivePeriod
func WithKeepAlivePeriod(keepAlivePeriod time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.KeepAlivePeriod = keepAlivePeriod
	}
}