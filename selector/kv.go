package selector

// Node defines the basic information for a service Node
type Node struct {
	Key string
	Value []byte
	weight int
}
