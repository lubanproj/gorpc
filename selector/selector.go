package selector

// Selector obtains a service node through service discovery and load balancing
type Selector interface {
	Select(string) (string, error)
}

type defaultSelector struct {

}

// Options defines Selector options
type Options struct {

}

type Option func(*Options)

func init() {
	RegisterSelector("default", DefaultSelector)
}

// DefaultSelector is a global default Selector
var DefaultSelector = &defaultSelector{}

var selectorMap = make(map[string]Selector)

// RegisterSelector supports business custom registered Selector
func RegisterSelector(name string, selector Selector) {
	if selectorMap == nil {
		selectorMap = make(map[string]Selector)
	}
	selectorMap[name] = selector
}

func (d *defaultSelector) Select(serviceName string) (string, error) {

	return "", nil
}

// GetSelector get a selector by a given selector name
func GetSelector(name string) Selector {
	if selector ,ok := selectorMap[name]; ok {
		return selector
	}
	return DefaultSelector
}

