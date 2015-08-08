package jarvis

import (
	"log"
	"net/rpc"
	"testing"
)

func TestServer(t *testing.T) {
	go NewServer()

	client, err := rpc.DialHTTP("tcp", "localhost:6060")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	var reply RunnerReply
	args := RunnerArgs{
		Files: []File{
			{
				Name:     "main.js",
				Contents: "console.log('Hello, World!')",
				Entry:    true,
			},
		},
		ProgramName: "Hello to the world!",
	}

	err = client.Call("Runner.Run", &args, &reply)
	if err != nil {
		log.Fatal("running: ", err)
	}

	log.Println(reply.Output)
}
