package connpool

type Options struct {
	initialCap int
	maxCap int
}

type Option func(*Options)

func WithInitialCap (initialCap int) Option {
	return func(o *Options) {
		o.initialCap = initialCap
	}
}

func WithMaxCap (maxCap int) Option {
	return func(o *Options) {
		o.maxCap = maxCap
	}
}
