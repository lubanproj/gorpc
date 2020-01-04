package gorpc

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// gorpc Server, 一个 Server 可以拥有一个或者多个 service
type Server struct {
	opts *ServerOptions
	services map[string]Service
}

func NewServer(opt ...ServerOption) *Server{

	s := &Server {
		opts : &ServerOptions{},
		services: make(map[string]Service),
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
		log.Fatalf("handlerType %v not match service : %v ", ht, st)
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

	for _, service := range s.services {
		go service.Serve(s.opts)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	<-ch

	s.Close()
}

func (s *Server) Close() {

}