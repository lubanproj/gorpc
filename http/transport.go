package http

import (
	"context"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/transport"
)

type httpServerTransport struct {
	http.Server
	opts *transport.ServerTransportOptions

	Router *httprouter.Router // router for httpServerTransport
}

// DefaultRouter uses httprouter as the default router
var DefaultRouter *httprouter.Router

func init() {
	// use httprouter for default router
	DefaultRouter = httprouter.New()
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


	s.Server.Handler = DefaultRouter
	go func() {
		if err = s.Server.Serve(lis); err != nil {
			log.Errorf("http serve error, %v", err)
		}
	}()

	return nil
}

// HandlerFunc is an adapter which allows the usage of an http handler
// request handle.
func HandleFunc(method, path string, handler func(http.ResponseWriter, *http.Request)) error {

	DefaultRouter.HandlerFunc(method, path, handler)

	return nil
}


