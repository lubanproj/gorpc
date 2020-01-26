package selector

type Resolver interface {
	Resolve(string) []*Node
}