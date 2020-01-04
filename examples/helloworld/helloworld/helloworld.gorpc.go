package helloworld

import (
	"context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/client"
	"github.com/lubanproj/gorpc/interceptor"
	math "math"
)

import (

)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

/* ************************************ Service Definition ************************************ */
type GreeterService interface {
	SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error)
}

var _Greeter_serviceDesc = &gorpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterService)(nil),
	Methods : []*gorpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    GreeterService_SayHello_Handler,
		},
	},
}

func GreeterService_SayHello_Handler(svr interface{},ctx context.Context, req interface{}, cep interceptor.ServerInterceptor) (interface{}, error) {

	if cep == nil {
		return svr.(GreeterService).SayHello(ctx, req.(*HelloRequest))
	}

	handler := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
		return svr.(GreeterService).SayHello(ctx, reqbody.(*HelloRequest))
	}

	return cep(ctx, req, handler)
}

func RegisterService(s *gorpc.Server, svr interface{}) {
	s.Register(_Greeter_serviceDesc, svr)
}

/* ************************************ Client Definition ************************************ */
type GreeterClientProxy interface {
	SayHello(ctx context.Context, req *HelloRequest, opts ...client.Option) (*HelloReply, error)
}

type GreeterClientProxyImpl struct {
	client client.Client
	opts   []client.Option
}

func NewGreeterClientProxy(opts ...client.Option) GreeterClientProxy {
	return &GreeterClientProxyImpl{client: client.DefaultClient, opts: opts}
}

func (c *GreeterClientProxyImpl) SayHello(ctx context.Context, req *HelloRequest,
	opts ...client.Option) (*HelloReply, error) {

	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)

	rsp := &HelloReply{}
	err := c.client.Invoke(ctx, req, rsp, callopts...)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
