package main

import (
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:6060")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	var reply HelloReply
	args := HelloArgs{"Harrison"}

	err = client.Call("HelloService.Goodbye", &args, &reply)
	if err != nil {
		log.Fatal("greeting: ", err)
	}

	log.Println(reply.Message)
}

type HelloArgs struct {
	Name string
}

type HelloReply struct {
	Message string
}
