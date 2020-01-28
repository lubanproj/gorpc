package gorpc

import (
	"github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/plugin"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// gorpc Server, 一个 Server 可以拥有一个或者多个 service
type Server struct {
	opts *ServerOptions
	services map[string]Service
	plugins []plugin.Plugin
}

func NewServer(opt ...ServerOption) *Server{

	s := &Server {
		opts : &ServerOptions{},
		services: make(map[string]Service),
	}
	for _, plugin := range s.plugins {
		s.plugins = append(s.plugins, plugin)
	}

	for _, o := range opt {
		o(s.opts)
	}

	return s
}

func (s *Server) Register(sd *ServiceDesc, svr interface{}) {
	if sd == nil || svr == nil {
		return
	}
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(svr)
	if !st.Implements(ht) {
		log.Fatal("handlerType %v not match service : %v ", ht, st)
	}

	ser := &service {
		svr : svr,
		serviceName : sd.ServiceName,
		handlers : make(map[string]Handler),
	}

	for _, method := range sd.Methods {
		ser.handlers[method.MethodName] = method.Handler
	}

	s.services[sd.ServiceName] = ser
}

func (s *Server) Serve() {

	// 加载所有插件
	for _, p := range s.plugins {
		if rp, ok := p.(plugin.ResolverPlugin); ok {

			var services []string
			for serviceName, _ := range s.services {
				services = append(services, serviceName)
			}

			pluginOptions := []plugin.Option {
				plugin.WithSelectorSvrAddr(s.opts.selectorSvrAddr),
				plugin.WithSvrAddr(s.opts.address),
				plugin.WithServices(services),
			}
			if err := rp.Init(pluginOptions ...); err != nil {
				log.Fatal("plugin init error, %v", err)
			}
		}
	}

	for _, service := range s.services {
		go service.Serve(s.opts)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	<-ch

	s.Close()
}

func (s *Server) Close() {

}