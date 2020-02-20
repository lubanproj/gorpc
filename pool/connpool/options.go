package connpool

import "time"

type Options struct {
	initialCap int   // initial capacity
	maxCap int    // max capacity
	idleTimeout time.Duration
	maxIdle int   // max idle connections
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


func WithMaxIdle (maxIdle int) Option {
	return func(o *Options) {
		o.maxIdle = maxIdle
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = idleTimeout
	}
}