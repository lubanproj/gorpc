package metadata

import "context"

type clientMD struct {}
type serverMD struct {}

type clientMetadata map[string][]byte

type serverMetadata map[string][]byte

// ClientMetadata creates a new context with key-value pairs attached.
func ClientMetadata(ctx context.Context) clientMetadata {
	if md, ok := ctx.Value(clientMD{}).(clientMetadata); ok {
		return md
	}
	md := make(map[string][]byte)
	WithClientMetadata(ctx, md)
	return md
}

// WithClientMetadata creates a new context with the specified metadata
func WithClientMetadata(ctx context.Context, metadata map[string][]byte) context.Context{
	return context.WithValue(ctx, clientMD{}, clientMetadata(metadata))
}

// ServerMetadata creates a new context with key-value pairs attached.
func ServerMetadata(ctx context.Context) serverMetadata {
	if md, ok := ctx.Value(serverMD{}).(serverMetadata); ok {
		return md
	}
	md := make(map[string][]byte)
	WithServerMetadata(ctx, md)
	return md
}

// WithServerMetadata creates a new context with the specified metadata
func WithServerMetadata(ctx context.Context, metadata map[string][]byte) context.Context{
	return context.WithValue(ctx, serverMD{}, serverMetadata(metadata))
}

