package jaeger

import (
	"context"
	"errors"
	"fmt"
	"github.com/lubanproj/gorpc/interceptor"
	"github.com/lubanproj/gorpc/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"
	"strings"
)

// Jaeger implements the opentracing specification
type Jaeger struct {
	opts *plugin.Options
}

const Name = "jaeger"
const JaegerClientName = "gorpc-client-jaeger"
const JaegerServerName = "gorpc-server-jaeger"

func init() {
	plugin.Register(Name, JaegerSvr)
}

// global jaeger objects for framework
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

// OpenTracingClientInterceptor packaging jaeger tracer as a client interceptor
func OpenTracingClientInterceptor(tracer opentracing.Tracer, spanName string) interceptor.ClientInterceptor {

	return func (ctx context.Context, req, rsp interface{}, ivk interceptor.Invoker) error {

		//var parentCtx opentracing.SpanContext
		//
		//if parent := opentracing.SpanFromContext(ctx); parent != nil {
		//	parentCtx = parent.Context()
		//}

		//clientSpan := tracer.StartSpan(spanName, ext.SpanKindRPCClient, opentracing.ChildOf(parentCtx))
		clientSpan := tracer.StartSpan(spanName, ext.SpanKindRPCClient)
		defer clientSpan.Finish()

		mdCarrier := &jaegerCarrier{}

		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdCarrier); err != nil {
			clientSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
		}

		clientSpan.LogFields(log.String("spanName", spanName))

		return ivk(ctx, req, rsp)

	}
}

// OpenTracingServerInterceptor packaging jaeger tracer as a server interceptor
func OpenTracingServerInterceptor(tracer opentracing.Tracer, spanName string) interceptor.ServerInterceptor {

	return func(ctx context.Context, req interface{}, handler interceptor.Handler) (interface{}, error) {

		mdCarrier := &jaegerCarrier{}

		spanContext, err := tracer.Extract(opentracing.HTTPHeaders, mdCarrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			return nil, errors.New(fmt.Sprintf("tracer extract error : %v", err))
		}
		serverSpan := tracer.StartSpan(spanName, ext.RPCServerOption(spanContext),ext.SpanKindRPCServer)
		defer serverSpan.Finish()

		ctx = opentracing.ContextWithSpan(ctx, serverSpan)

		serverSpan.LogFields(log.String("spanName", spanName))

		return handler(ctx, req)
	}

}

// Init implements the initialization of the jaeger configuration when the framework is loaded
func Init(tracingSvrAddr string, opts ... plugin.Option) (opentracing.Tracer, error) {
	return initJaeger(tracingSvrAddr, JaegerClientName, opts ...)
}

func (j *Jaeger) Init(opts ...plugin.Option) (opentracing.Tracer, error) {

	for _, o := range opts {
		o(j.opts)
	}

	if j.opts.TracingSvrAddr == "" {
		return nil, errors.New("jaeger init error, traingSvrAddr is empty")
	}

	return initJaeger(j.opts.TracingSvrAddr, JaegerServerName, opts ...)

}

func initJaeger(tracingSvrAddr string, jaegerServiceName string, opts ... plugin.Option) (opentracing.Tracer, error) {
	cfg := &config.Configuration{
		Sampler : &config.SamplerConfig{
			Type : "const",  // Fixed sampling
			Param : 1,       // 1= full sampling, 0= no sampling
		},
		Reporter : &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: tracingSvrAddr,
		},
		ServiceName : jaegerServiceName,
	}

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, err
}
