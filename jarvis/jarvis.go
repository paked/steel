package jarvis

import (
	"bytes"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
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
	Name     string // hello_world
	Type     string // js
	Contents string // console.log('Hello, World!');
}

type RunnerArgs struct {
	ProgramName string
	Main        File
	Resources   []File
}

type RunnerReply struct {
	OK     bool
	Output string
}

type Runner struct {
}

func (r *Runner) Run(args *RunnerArgs, reply *RunnerReply) error {
	if args.Main == (File{}) {
		return errors.New("No main file!")
	}

	switch args.Main.Type {
	case "js":
		f, err := os.Create(os.TempDir() + args.Main.Name + ".js")
		if err != nil {
			return err
		}

		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		_, err = f.WriteString(args.Main.Contents)
		if err != nil {
			return err
		}

		cmd := exec.Command("node", f.Name())
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			return err
		}

		*reply = RunnerReply{
			OK:     true,
			Output: out.String(),
		}
	default:
		return errors.New("We do not support that programming language!")
	}

	return nil
}
