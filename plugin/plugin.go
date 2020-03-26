package plugin

import "github.com/opentracing/opentracing-go"

// Plugin defines the standard for all plug-ins
type Plugin interface {

}

// ResolverPlugin defines the standard for all server discovery plug-ins
type ResolverPlugin interface {
	Init(...Option) error
}

// TracingPlugin defines the standard for all tracing plug-ins
type TracingPlugin interface {
	Init(...Option) (opentracing.Tracer, error)
}

// PluginMap defines a global plug-in map
var PluginMap = make(map[string]Plugin)

// Register opens an entry point for all plug-ins to register
func Register(name string, plugin Plugin) {
	if PluginMap == nil {
		PluginMap = make(map[string]Plugin)
	}
	PluginMap[name] = plugin
}

// Options for all plug-ins
type Options struct {
	SvrAddr string     // server address
	Services []string   // service arrays
	SelectorSvrAddr string  // server discovery address ，e.g. consul server address
	TracingSvrAddr string   // tracing server address，e.g. jaeger server address
}

// Option provides operations on Options
type Option func(*Options)

// WithSvrAddr allows you to set SvrAddr of Options
func WithSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SvrAddr = addr
	}
}

// WithSvrAddr allows you to set Services of Options
func WithServices(services []string) Option {
	return func(o *Options) {
		o.Services = services
	}
}

// WithSvrAddr allows you to set SelectorSvrAddr of Options
func WithSelectorSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SelectorSvrAddr = addr
	}
}

// WithSvrAddr allows you to set TracingSvrAddr of Options
func WithTracingSvrAddr(addr string) Option {
	return func(o *Options) {
		o.TracingSvrAddr = addr
	}
}




