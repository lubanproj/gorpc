package jaeger

import (
	"context"
	"github.com/lubanproj/gorpc/interceptor"
	gorpclog "github.com/lubanproj/gorpc/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)


type jaegerCarrier map[string][]byte

func (m jaegerCarrier) Set(key, val string) {
	key = strings.ToLower(key)
	m[key] = []byte(val)
}

func (m jaegerCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range m {
		handler(k, string(v))
	}
	return nil
}


func OpenTracingClientInterceptor(tracer opentracing.Tracer, servicePath string) interceptor.ClientInterceptor {

	return func (ctx context.Context, req, rsp interface{}, ivk interceptor.Invoker) error {

		var parentCtx opentracing.SpanContext

		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}

		clientSpan := tracer.StartSpan(servicePath, ext.SpanKindRPCClient, opentracing.ChildOf(parentCtx))
		defer clientSpan.Finish()

		mdCarrier := &jaegerCarrier{}

		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdCarrier); err != nil {
			clientSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
		}

		return ivk(ctx, req, rsp)

	}
}

func OpenTracingServerInterceptor(tracer opentracing.Tracer, servicePath string) interceptor.ServerInterceptor {

	return func(ctx context.Context, req interface{}, handler interceptor.Handler) (interface{}, error) {

		mdCarrier := &jaegerCarrier{}

		spanContext, err := tracer.Extract(opentracing.HTTPHeaders, mdCarrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			gorpclog.Error("Tracer.Extract() failed, %v", err)
		}
		serverSpan := tracer.StartSpan(servicePath, ext.RPCServerOption(spanContext),ext.SpanKindRPCServer)
		defer serverSpan.Finish()

		ctx = opentracing.ContextWithSpan(ctx, serverSpan)

		return handler(ctx, req)
	}

}