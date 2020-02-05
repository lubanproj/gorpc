package plugin

import "github.com/opentracing/opentracing-go"

// 插件
type Plugin interface {

}

// 服务发现类插件
type ResolverPlugin interface {
	Init(...Option) error
}

// Tracing 类插件
type TracingPlugin interface {
	Init(...Option) (opentracing.Tracer, error)
}

var PluginMap = make(map[string]Plugin)

func Register(name string, plugin Plugin) {
	if PluginMap == nil {
		PluginMap = make(map[string]Plugin)
	}
	PluginMap[name] = plugin
}

type Options struct {
	SvrAddr string     // server 地址
	Services []string   // 服务名数组
	SelectorSvrAddr string  // 服务发现集群地址 ，例如 consul server 地址
	TracingSvrAddr string   // tracing server 地址，例如 jaeger server 地址
}

type Option func(*Options)

func WithSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SvrAddr = addr
	}
}

func WithServices(services []string) Option {
	return func(o *Options) {
		o.Services = services
	}
}

func WithSelectorSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SelectorSvrAddr = addr
	}
}

func WithTracingSvrAddr(addr string) Option {
	return func(o *Options) {
		o.TracingSvrAddr = addr
	}
}




