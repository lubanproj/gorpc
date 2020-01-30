package main

import (
	"context"
	"fmt"
	"github.com/lubanproj/gorpc/client"
	"github.com/lubanproj/gorpc/examples/helloworld/helloworld"
	"time"
)

func main() {
	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(2000 * time.Millisecond),
	}
	c := client.DefaultClient
	req := &helloworld.HelloRequest{
		Msg: "hello",
	}
	rsp := &helloworld.HelloReply{}
	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)
	fmt.Println(rsp.Msg, err)
}
