package interceptor

import "context"

type ServerInterceptor func(ctx context.Context, req interface{}, handler Handler) (interface{}, error)

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type ClientInterceptor func(ctx context.Context, req, rsp interface{}, ivk Invoker) error

type Invoker func(ctx context.Context, req, rsp interface{}) error

func ClientIntercept(ctx context.Context, req, rsp interface{}, ceps []ClientInterceptor, ivk Invoker) error {

	if len(ceps) == 0 {
		return ivk(ctx, req, rsp)
	}

	return ceps[0](ctx, req, rsp, getInvoker(0, ceps, ivk))
}

func getInvoker(cur int, ceps []ClientInterceptor, ivk Invoker) Invoker {
	if cur == len(ceps) - 1 {
		return ivk
	}

	return func(ctx context.Context, req , rsp interface{} ) error {
		return ceps[cur+1](ctx, req, rsp, getInvoker(cur+1, ceps, ivk))
	}
}

func ServerIntercept(ctx context.Context, req interface{}, ceps []ServerInterceptor, handler Handler) (interface{}, error) {

	if len(ceps) == 0 {
		return handler(ctx, req)
	}

	return ceps[0](ctx, req, getHandler(0, ceps, handler))
}

func getHandler(cur int, ceps []ServerInterceptor, handler Handler) Handler {
	if cur == len(ceps) - 1 {
		return handler
	}

	return func(ctx context.Context, req interface{} ) (interface{}, error) {
		return ceps[cur+1](ctx, req, getHandler(cur+1, ceps, handler))
	}
}