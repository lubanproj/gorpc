package main

import (
	"context"
	"fmt"
	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/examples/helloworld2/helloworld"
	"time"
)

type greeterService struct{}

func (g *greeterService) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println("recv Msg : ", req.Msg)
	rsp := &helloworld.HelloReply{
		Msg: "hello, " + req.Msg ,
	}
	return rsp, nil
}


func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithProtocol("proto"),
		gorpc.WithTimeout(time.Millisecond * 2000),
	}
	s := gorpc.NewServer(opts ...)
	helloworld.RegisterService(s, &greeterService{})
	s.Serve()
}