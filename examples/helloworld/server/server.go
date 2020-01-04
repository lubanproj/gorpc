package main

import (
	"context"
	"fmt"
	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/examples/helloworld/helloworld"
)

type greeterService struct{}

func (g *greeterService) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println("recv req : %v", req)
	rsp := &helloworld.HelloReply{
		Msg: "hello, " + req.Msg ,
	}
	return rsp, nil
}


func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithTarget("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithProtocol("proto"),
	}
	s := gorpc.NewServer(opts ...)
	helloworld.RegisterService(s, &greeterService{})
	s.Serve()
}