package selector

type Selector interface {
	Select() string
}

type defaultSelector struct {

}

func (d *defaultSelector) Select(serviceName string) string {

	return serviceName
}
