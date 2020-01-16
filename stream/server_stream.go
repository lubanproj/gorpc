package stream

import "context"

type ServerStream struct {
	ctx context.Context
	Method string // 方法名
	RetCode uint32 // 返回码 0—成功 非0-失败
	RetMsg  string  // 返回信息 OK-成功，失败返回具体信息
}

const ServerStreamKey = StreamContextKey("GORPC_SERVER_STREAM")

func GetServerStream(ctx context.Context) *ServerStream {
	v := ctx.Value(ServerStreamKey)
	if v == nil {
		cs := &ClientStream{}
		context.WithValue(ctx, ServerStreamKey, cs)
	}
	return v.(*ServerStream)
}

func (ss *ServerStream) WithMethod(method string) *ServerStream {
	ss.Method = method
	return ss
}

func (ss *ServerStream) Clone() *ServerStream {
	return &ServerStream{
		Method : ss.Method,
	}
}


func NewServerStream(ctx context.Context) (context.Context, *ServerStream) {
	var ss *ServerStream
	v := ctx.Value(ServerStreamKey)
	if v != nil {
		ss = v.(*ServerStream)
	} else {
		ss = &ServerStream{
			ctx : ctx,
		}
	}
	valueCtx := context.WithValue(ctx, ServerStreamKey, ss)
	return valueCtx, ss
}
