package interceptor

import "context"

type Interceptor func(ctx context.Context, req interface{}, rsp interface{}, ivk Invoker ) error

type Invoker func (ctx context.Context, req interface{}, rsp interface{} ) error


func Intercept(ctx context.Context, req interface{}, rsp interface{},ceps []Interceptor, ivk Invoker ) error {

	if len(ceps) == 0 {
		return ivk(ctx, req, rsp)
	}

	return ceps[0](ctx, req, rsp, getInvoker(ctx, req, rsp, 0, ceps, ivk))

}

func getInvoker(ctx context.Context, req interface{}, rsp interface{}, cur int, ceps []Interceptor, ivk Invoker) Invoker {
	if cur == len(ceps) - 1 {
		return ivk
	}

	return func(ctx context.Context, req interface{}, rsp interface{} ) error {
		return ceps[cur+1](ctx, req, rsp, getInvoker(ctx, req, rsp, cur+1, ceps, ivk))
	}
}
