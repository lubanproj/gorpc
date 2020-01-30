package gorpc

import (
	"context"
	"fmt"
	"github.com/lubanproj/gorpc/interceptor"
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

	for pluginName, plugin := range plugin.PluginMap {
		if !containPlugin(pluginName, s.opts.pluginNames) {
			continue
		}
		s.plugins = append(s.plugins, plugin)
	}

	for _, o := range opt {
		o(s.opts)
	}

	return s
}

func containPlugin(pluginName string, plugins []string) bool {
	for _, plugin := range plugins {
		if pluginName == plugin {
			return true
		}
	}
	return false
}

type emptyInterface interface{}

func (s *Server) RegisterService(serviceName string, svr interface{}) error {

	svrType := reflect.TypeOf(svr)
	svrValue := reflect.ValueOf(svr)

	sd := &ServiceDesc{
		ServiceName: serviceName,
		// 这里为了和代码生成兼容
		HandlerType : (*emptyInterface)(nil),
		Svr : svr,
	}

	methods, err := getServiceMethods(svrType, svrValue)
	if err != nil {
		return err
	}

	sd.Methods = methods

	s.Register(sd, svr)

	return nil
}

func getServiceMethods(serviceType reflect.Type, serviceValue reflect.Value) ([]*MethodDesc, error) {

	var methods []*MethodDesc

	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)

		if err := checkMethod(method.Type); err != nil {
			return nil, err
		}

		methodHandler := func (svr interface{},ctx context.Context, dec func(interface{}) error, cep interceptor.ServerInterceptor) (interface{}, error) {

			reqType := method.Type.In(2)

			// 判断类型
			req := reflect.New(reqType.Elem()).Interface()

			if err := dec(req); err != nil {
				return nil, err
			}

			if cep == nil {
				values := method.Func.Call([]reflect.Value{serviceValue,reflect.ValueOf(ctx),reflect.ValueOf(req)})
				// 判断错误
				return values[0].Interface(), nil
			}

			handler := func(ctx context.Context, reqbody interface{}) (interface{}, error) {

				// 执行反射
				values := method.Func.Call([]reflect.Value{serviceValue,reflect.ValueOf(ctx),reflect.ValueOf(req)})

				// 判断错误
				return values[0].Interface(), nil
			}

			return cep(ctx, req, handler)
		}

		methods = append(methods, &MethodDesc{
			MethodName: method.Name,
			Handler: methodHandler,
		})
	}

	return methods , nil
}

func checkMethod(method reflect.Type) error {

	// 参数个数 >= 2 , 这里需要加上自身
	if method.NumIn() < 3 {
		return fmt.Errorf("method %s invalid, the number of params < 2", method.Name())
	}

	// 返回值个数为 2
	if method.NumOut() != 2 {
		return fmt.Errorf("method %s invalid, the number of return values != 2", method.Name())
	}

	// 第一个参数必须是 context
	ctxType := method.In(1)
	var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	if !ctxType.Implements(contextType) {
		return fmt.Errorf("method %s invalid, first param is not context", method.Name())
	}

	// 第二个参数必须是指针
	argType := method.In(2)
	if argType.Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, req type is not a pointer", method.Name())
	}

	// 第一个返回值必须是指针
	replyType := method.Out(0)
	if replyType.Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, reply type is not a pointer", method.Name())
	}

	// 第二个返回值必须是 error
	errType := method.Out(1)
	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	if !errType.Implements(errorType) {
		return fmt.Errorf("method %s invalid, returns %s , not error", method.Name(), errType.Name())
	}

	return nil
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