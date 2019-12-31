package main

import (
	"github.com/lubanproj/gorpc"
)

func main() {
	opts := []gorpc.ServiceOption{
		gorpc.WithTarget("127.0.0.1:8000"),
		gorpc.WithNetwork("tcp"),
	}
	s := gorpc.NewServer(opts ...)
	s.Serve()
}