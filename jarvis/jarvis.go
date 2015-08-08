package jarvis

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func NewServer() {
	runner := &Runner{}
	rpc.Register(runner)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":6060")
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	http.Serve(l, nil)
}

type File struct {
	Name     string
	Contents string
	Entry    bool
}

type RunnerArgs struct {
	ProgramName string
	Files       []File
}

type RunnerReply struct {
	OK     bool
	Output string
}

type Runner struct {
}

func (r *Runner) Run(args *RunnerArgs, reply *RunnerReply) error {
	if len(args.Files) == 0 {
		return errors.New("No files submitted")
	}

	*reply = RunnerReply{
		OK:     true,
		Output: "Hello, World!",
	}

	return nil
}
