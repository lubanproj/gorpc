package gorpc

import (
	"context"
	"github.com/luban_proj/gorpc/log"
)

// Service 定义了某个具体服务的通用实现接口
type Service interface {
	Register()
	Serve()
	Close()
}


type service struct{
	ctx context.Context  		// 每一个 service 一个上下文进行管理
	cancel context.CancelFunc   // context 的控制器
	serviceName string   		// 服务名
	handlers map[string]Handler
	opts ServiceOptions  		// 参数选项

}

type Handler interface {
	Handle(context.Context, []byte) ([]byte, error)
}


func (s *service) Register(handlerName string, handler Handler) {
	s.handlers[handlerName] = handler
}

func (s *service) Serve() {

	if err := s.opts.transport.ListenAndServe(s.ctx, s.opts.transportOptions ...); err != nil {
		log.Error("%s serve error, %v", s.serviceName, err)
		return
	}

	<- s.ctx.Done()
}


func (s *service) Close() {

}

