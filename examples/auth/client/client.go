package main

import (
	"context"
	"fmt"
	"github.com/lubanproj/gorpc/auth"
	"time"

	"github.com/lubanproj/gorpc/client"
	"github.com/lubanproj/gorpc/testdata"
)

func main() {
	opts := []client.Option {
		client.WithTarget("127.0.0.1:8003"),
		client.WithNetwork("tcp"),
		client.WithTimeout(2000000 * time.Millisecond),
		client.WithSerializationType("msgpack"),
		client.WithPerRPCAuth(auth.NewOAuth2ByToken("testToken")),
	}
	c := client.DefaultClient
	req := &testdata.HelloRequest{
		Msg: "hello",
	}
	rsp := &testdata.HelloReply{}
	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, rsp, opts ...)
	fmt.Println(rsp.Msg, err)
}
