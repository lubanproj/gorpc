package selector

type Resolver interface {
	Resolve() []*Node
}