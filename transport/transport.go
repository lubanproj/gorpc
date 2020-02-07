// Network communication layer, responsible for the bottom layer of network communication,
// mainly including tcp && udp two protocol implementation
package transport

import "context"

type ServerTransport interface {
	// monitoring and processing of requests
	ListenAndServe(context.Context, ...ServerTransportOption) error
}

type ClientTransport interface {
	// send requests
	Send(context.Context, []byte, ...ClientTransportOption) ([]byte, error)
}


