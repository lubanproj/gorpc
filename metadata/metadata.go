package metadata

import "context"

type clientMD struct {}
type serverMD struct {}

type clientMetadata map[string][]byte

type serverMetadata map[string][]byte

func ClientMetadata(ctx context.Context) clientMetadata {
	if md, ok := ctx.Value(clientMD{}).(clientMetadata); ok {
		return md
	}
	md := make(map[string][]byte)
	WithClientMetadata(ctx, md)
	return md
}

func WithClientMetadata(ctx context.Context, metadata map[string][]byte) context.Context{
	return context.WithValue(ctx, clientMD{}, metadata)
}


func ServerMetadata(ctx context.Context) serverMetadata {
	if md, ok := ctx.Value(serverMD{}).(serverMetadata); ok {
		return md
	}
	md := make(map[string][]byte)
	WithServerMetadata(ctx, md)
	return md
}


func WithServerMetadata(ctx context.Context, serverMetadata map[string][]byte) context.Context{
	return context.WithValue(ctx, serverMD{}, serverMetadata)
}

