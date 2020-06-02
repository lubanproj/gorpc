package gorpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/metadata"
	"github.com/lubanproj/gorpc/protocol"
	"github.com/lubanproj/gorpc/transport"
	"github.com/lubanproj/gorpc/utils"

	"github.com/golang/protobuf/proto"
)

// Service defines a generic implementation interface for a specific Service
type Service interface {
	Register(string, Handler)
	Serve(*ServerOptions)
	Close()
}

type service struct {
	svr         interface{}        // server
	ctx         context.Context    // Each service is managed in one context
	cancel      context.CancelFunc // controller of context
	serviceName string             // service name
	handlers    map[string]Handler
	opts        *ServerOptions // parameter options

	closing bool // whether the service is closing
}

// ServiceDesc is a detailed description of a service
type ServiceDesc struct {
	Svr         interface{}
	ServiceName string
	Methods     []*MethodDesc
	HandlerType interface{}
}

// MethodDesc is a detailed description of a method
type MethodDesc struct {
	MethodName string
	Handler    Handler
}

// Handler is the handler of a method
type Handler func(context.Context, interface{}, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)

func (s *service) Register(handlerName string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handlerName] = handler
}

func (s *service) Serve(opts *ServerOptions) {

	s.opts = opts

	transportOpts := []transport.ServerTransportOption{
		transport.WithServerAddress(s.opts.address),
		transport.WithServerNetwork(s.opts.network),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.timeout),
		transport.WithSerializationType(s.opts.serializationType),
		transport.WithProtocol(s.opts.protocol),
	}

	serverTransport := transport.GetServerTransport(s.opts.protocol)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := serverTransport.ListenAndServe(s.ctx, transportOpts...); err != nil {
		log.Errorf("%s serve error, %v", s.opts.network, err)
		return
	}

	fmt.Printf("%s service serving at %s ... \n", s.opts.protocol, s.opts.address)

	<-s.ctx.Done()
}

func (s *service) Close() {
	s.closing = true
	if s.cancel != nil {
		s.cancel()
	}
	fmt.Println("service closing ...")
}

func (s *service) Handle(ctx context.Context, reqbuf []byte) ([]byte, error) {

	// parse protocol header
	request := &protocol.Request{}
	if err := proto.Unmarshal(reqbuf, request); err != nil {
		return nil, err
	}

	ctx = metadata.WithServerMetadata(ctx, request.Metadata)

	serverSerialization := codec.GetSerialization(s.opts.serializationType)

	dec := func(req interface{}) error {

		if err := serverSerialization.Unmarshal(request.Payload, req); err != nil {
			return err
		}
		return nil
	}

	if s.opts.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.opts.timeout)
		defer cancel()
	}

	_, method, err := utils.ParseServicePath(string(request.ServicePath))
	if err != nil {
		return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	}

	handler := s.handlers[method]
	if handler == nil {
		return nil, errors.New("handlers is nil")
	}

	rsp, err := handler(ctx, s.svr, dec, s.opts.interceptors)
	if err != nil {
		return nil, err
	}

	rspbuf, err := serverSerialization.Marshal(rsp)
	if err != nil {
		return nil, err
	}

	return rspbuf, nil
}
