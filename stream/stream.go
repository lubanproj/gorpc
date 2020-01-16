package stream

type StreamContextKey string

type Stream interface {
	Clone() Stream
}