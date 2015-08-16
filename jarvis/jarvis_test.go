package jarvis

import (
	"log"
	"net/rpc"
	"testing"
)

func TestServer(t *testing.T) {
	NewServer()

	client, err := rpc.DialHTTP("tcp", "localhost:6060")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	var reply RunnerReply
	args := RunnerArgs{
		Main: File{
			Name:     "main",
			Type:     "js",
			Contents: "console.log('Hello, World!')",
		},
		ProgramName: "Hello to the world!",
	}

	err = client.Call("Runner.Run", &args, &reply)
	if err != nil {
		log.Fatal("running: ", err)
	}

	log.Println(reply.Output)
}
