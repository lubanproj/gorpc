package main

import (
	"github.com/lubanproj/gorpc"
	"github.com/lubanproj/gorpc/examples/helloworld/helloworld"
	"github.com/lubanproj/gorpc/log"
	"time"
)


func main() {
	opts := []gorpc.ServerOption{
		gorpc.WithAddress("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
		gorpc.WithProtocol("proto"),
		gorpc.WithTimeout(time.Minute * 60),
	}
	s := gorpc.NewServer(opts ...)
	if err := s.RegisterService("/helloworld.Greeter", new(helloworld.Service)); err != nil {
		log.Fatal("register service error, %v", err)
	}
	s.Serve()
}