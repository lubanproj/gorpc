package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lubanproj/gorpc/client"
	"github.com/lubanproj/gorpc/plugin/jaeger"
	"github.com/lubanproj/gorpc/testdata"

)

func main() {

	tracer, err := jaeger.Init("localhost:6831")
	if err != nil {
		panic(err)
	}

	opts := []client.Option {
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(2000 * time.Millisecond),
		client.WithInterceptor(jaeger.OpenTracingClientInterceptor(tracer, "/helloworld.Greeter/SayHello")),
	}
	c := client.DefaultClient
	req := &testdata.HelloRequest{
		Msg: "hello",
	}
	rsp := &testdata.HelloReply{}

	for i:= 1; i< 200; i ++ {
		err = c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)
		fmt.Println(rsp.Msg, err)
		time.Sleep(100 * time.Millisecond)
	}

}
