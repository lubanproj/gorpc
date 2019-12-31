package interceptor

import "context"

type Interceptor func(ctx context.Context, req interface{}, ivk Invoker ) (interface{}, error)

type Invoker func (ctx context.Context, req interface{}) (interface{}, error)


func Intercept(ctx context.Context, req interface{},ceps []Interceptor, ivk Invoker ) (interface{},error) {

	if len(ceps) == 0 {
		return ivk(ctx, req)
	}

	return ceps[0](ctx, req, getInvoker(ctx, req, 0, ceps, ivk))

}

func getInvoker(ctx context.Context, req interface{}, cur int, ceps []Interceptor, ivk Invoker) Invoker {
	if cur == len(ceps) - 1 {
		return ivk
	}

	return func(ctx context.Context, req interface{} ) (interface{}, error) {
		return ceps[cur+1](ctx, req, getInvoker(ctx, req, cur+1, ceps, ivk))
	}
}
