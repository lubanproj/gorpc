package jaeger

import (
	"context"
	"errors"
	"github.com/lubanproj/gorpc/interceptor"
	gorpclog "github.com/lubanproj/gorpc/log"
	"github.com/lubanproj/gorpc/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"
	"strings"
)


type Jaeger struct {
	opts *plugin.Options
}

const Name = "jaeger"

func init() {
	plugin.Register(Name, JaegerSvr)
}

var JaegerSvr = &Jaeger {
	opts : &plugin.Options{},
}

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


func Init(tracingSvrAddr string, opts ... plugin.Option) error {
	cfg := &config.Configuration{}

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)
	return err
}

func (c *Jaeger) Init(opts ...plugin.Option) error {

	for _, o := range opts {
		o(c.opts)
	}

	if c.opts.TracingSvrAddr == "" {
		return errors.New("jaeger init error, traingSvrAddr is empty")
	}

	cfg := &config.Configuration{}

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	return nil

}