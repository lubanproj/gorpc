package auth

import (
	"context"
	"net"
)

// TransportAuth defines a common interface for client and server handshakes
type TransportAuth interface {

	// ClientHandshake defines a common interface for client handshakes
	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)

	// ServerHandshake defines a common interface for server handshakes
	ServerHandshake(conn net.Conn) (net.Conn, AuthInfo, error)

}

// PerRPCAuth defines a common interface for single RPC call authentication
type PerRPCAuth interface {

	// GetMetadata fetch custom metadata from the context
	GetMetadata(ctx context.Context, uri ... string) (map[string]string, error)

}

// AuthInfo defines the protocol type for authentication
type AuthInfo interface {
	AuthType() string
}