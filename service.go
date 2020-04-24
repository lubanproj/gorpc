package gorpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/codes"
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/protocol"
	"github.com/lubanproj/gorpc/stream"
	"github.com/lubanproj/gorpc/transport"
	"github.com/lubanproj/gorpc/utils"
)

// Service 定义了某个具体服务的通用实现接口
type Service interface {
	Register(string, Handler)
	Serve(*ServerOptions)
	Close()
}

type service struct{
	svr interface{}  			// server
	ctx context.Context  		// 每一个 service 一个上下文进行管理
	cancel context.CancelFunc   // context 的控制器
	serviceName string   		// 服务名
	handlers map[string]Handler
	opts *ServerOptions  		// 参数选项

	closing bool    // whether the service is closing
}

// ServiceDesc is a detailed description of a service
type ServiceDesc struct {
	Svr interface{}
	ServiceName string
	Methods []*MethodDesc
	HandlerType interface{}
}

// MethodDesc is a detailed description of a method
type MethodDesc struct {
	MethodName string
	Handler Handler
}

// Handler is the handler of a method
type Handler func (interface{}, context.Context, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)

func (s *service) Register(handlerName string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handlerName] = handler
}

func (s *service) Serve(opts *ServerOptions) {

	s.opts = opts

	transportOpts := []transport.ServerTransportOption {
		transport.WithServerAddress(s.opts.address),
		transport.WithServerNetwork(s.opts.network),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.timeout),
		transport.WithSerialization(s.opts.serializationType),
	}

	serverTransport := transport.GetServerTransport(s.opts.protocol)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := serverTransport.ListenAndServe(s.ctx, transportOpts ...); err != nil {
		log.Errorf("%s serve error, %v", s.opts.network, err)
		return
	}

	fmt.Printf("%s service serving at %s ... \n",s.opts.protocol, s.opts.address)

	<- s.ctx.Done()
}

func (s *service) Close() {
	s.closing = true
	s.cancel()
	fmt.Println("service closing ...")
}


func (s *service) Handle (ctx context.Context, frame []byte) ([]byte, error) {

	if len(frame) == 0 {
		return nil, errors.New("req is nil")
	}

	// 将 reqbuf 解析成 req interface {}
	serverCodec := codec.GetCodec(s.opts.protocol)
	reqbuf, err := serverCodec.Decode(frame)
	if err != nil {
		return nil, err
	}

	// parse protocol header
	request := &protocol.Request{}
	if err = proto.Unmarshal(reqbuf, request); err != nil {
		return nil, err
	}

	serverSerialization := codec.GetSerialization(s.opts.serializationType)

	dec := func(req interface {}) error {

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

	_, method , err := utils.ParseServicePath(string(request.ServicePath))
	if err != nil {
		return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	}

	handler := s.handlers[method]
	if handler == nil {
		return nil, errors.New("handlers is nil")
	}

	rsp, err := handler(s.svr, ctx, dec, s.opts.interceptors)
	if err != nil {
		return nil, err
	}

	rspbuf, err := serverSerialization.Marshal(rsp)
	if err != nil {
		return nil, err
	}

	response := addRspHeader(ctx, rspbuf)
	rspPb, err := proto.Marshal(response)
	if err != nil {
		return nil, err
	}

	rspbody, err := serverCodec.Encode(rspPb)
	if err != nil {
		return nil, err
	}

	return rspbody, nil
}


func addRspHeader(ctx context.Context, payload []byte) *protocol.Response {
	serverStream := stream.GetServerStream(ctx)
	response := &protocol.Response{
		Payload: payload,
		RetCode: serverStream.RetCode,
		RetMsg: serverStream.RetMsg,
	}

	return response
}
