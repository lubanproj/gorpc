package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lubanproj/gorpc/client"
	"github.com/lubanproj/gorpc/plugin/consul"
	"github.com/lubanproj/gorpc/testdata"
)

func main() {
	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(2000 * time.Millisecond),
		client.WithSelectorName(consul.Name),
	}
	c := client.DefaultClient
	req := &testdata.HelloRequest{
		Msg: "hello",
	}
	rsp := &testdata.HelloReply{}

	consul.Init("localhost:8500")
	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)
	fmt.Println(rsp.Msg, err)
}
