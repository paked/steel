package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	hello := &HelloService{}
	rpc.Register(hello)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":6060")
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	http.Serve(l, nil)
}

type HelloArgs struct {
	Name string
}

type HelloReply struct {
	Message string
}

type HelloService struct {
}

func (h *HelloService) Hello(args *HelloArgs, reply *HelloReply) error {
	reply.Message = fmt.Sprintf("Hello, %s!", args.Name)

	return nil
}

func (h *HelloService) Goodbye(args *HelloArgs, reply *HelloReply) error {
	reply.Message = fmt.Sprintf("Goodbye, %s!", args.Name)

	return nil
}
