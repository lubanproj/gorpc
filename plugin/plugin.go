package plugin

// 插件
type Plugin interface {

}

// 服务发现类插件
type ResolverPlugin interface {
	Init(...Option) error
}

var PluginMap = make(map[string]Plugin)

func Register(name string, plugin Plugin) {
	if PluginMap == nil {
		PluginMap = make(map[string]Plugin)
	}
	PluginMap[name] = plugin
}

type Options struct {
	SelectorSvrAddr string  // 服务发现集群地址 ，例如 consul server 地址
	SvrAddr string     // server 地址
	services []string   // 服务名数组
}

type Option func(*Options)

func WithSelectorSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SelectorSvrAddr = addr
	}
}

func WithSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SvrAddr = addr
	}
}

func WithServices(services []string) Option {
	return func(o *Options) {
		o.services = services
	}
}




