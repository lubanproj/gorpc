package http

import (
	"context"
	"github.com/lubanproj/gorpc/transport"
	"net"
	"net/http"
)

type httpServerTransport struct {
	http.Server
	opts *transport.ServerTransportOptions
}

// The default httpServerTransport
var DefaultHttpServerTransport = NewHttpServerTransport()

// Use the singleton pattern to create a server transport
var NewHttpServerTransport = func() *httpServerTransport {
	return &httpServerTransport{
		opts : &transport.ServerTransportOptions{},
	}
}

func init() {
	transport.RegisterServerTransport("http", DefaultHttpServerTransport)
}

func (s *httpServerTransport) ListenAndServe(ctx context.Context, opts ...transport.ServerTransportOption) error {
	for _, o := range opts {
		o(s.opts)
	}

	lis, err := net.Listen(s.opts.Network, s.opts.Address)
	if err != nil {
		return err
	}
	return s.Server.Serve(lis)
}



