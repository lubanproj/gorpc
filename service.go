package gorpc

import (
	"context"
	"errors"
	"github.com/lubanproj/gorpc/codec"
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/transport"
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
	ceps interceptor.ServerInterceptor
}

type ServiceDesc struct {
	Svr interface{}
	ServiceName string
	Methods []*MethodDesc
	HandlerType interface{}
}

type MethodDesc struct {
	MethodName string
	Handler Handler
}

type Handler func (interface{}, context.Context, func(interface{}) error, interceptor.ServerInterceptor) (interface{}, error)

func (s *service) Register(handlerName string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handlerName] = handler
}

func (s *service) Serve(opts *ServerOptions) {
	// TODO 思考下除了 Server 和 Service 的 Options 如何处理
	s.opts = opts

	transportOpts := []transport.ServerTransportOption {
		transport.WithServerAddress(s.opts.address),
		transport.WithServerNetwork(s.opts.network),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.timeout),
	}

	serverTransport := transport.GetServerTransport("default")

	if err := serverTransport.ListenAndServe(s.ctx, transportOpts ...); err != nil {
		log.Error("%s serve error, %v", s.serviceName, err)
		return
	}

	<- s.ctx.Done()
}

func (s *service) Close() {

}


func (s *service) Handle (ctx context.Context, reqbuf []byte) ([]byte, error) {

	if len(reqbuf) == 0 {
		return nil, errors.New("req is nil")
	}

	handler := s.handlers[s.opts.method]
	if handler == nil {
		return nil, errors.New("handlers is nil")
	}

	// 将 reqbuf 解析成 req interface {}
	serverCodec := codec.GetCodec(s.opts.protocol)
	serverSerialization := codec.GetSerialization(s.opts.protocol)

	dec := func(req interface {}) error {
		reqbody, err := serverCodec.Decode(reqbuf)
		if err != nil {
			return err
		}
		if err = serverSerialization.Unmarshal(reqbody, req); err != nil {
			return err
		}
		return nil
	}

	rsp, err := handler(s.svr, ctx, dec, s.ceps)

	rspbuf, err := serverSerialization.Marshal(rsp)
	if err != nil {
		return nil, err
	}

	rspbody, err := serverCodec.Encode(rspbuf)
	if err != nil {
		return nil, err
	}

	return rspbody, nil
}

