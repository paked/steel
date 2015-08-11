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
	OK        bool
	Output    string
	ErrOutput string
}

type Runner struct {
}

func (r *Runner) Run(args *RunnerArgs, reply *RunnerReply) error {
	if args.Main == (File{}) {
		return errors.New("No main file!")
	}

	switch args.Main.Type {
	case "js":
		var files []*os.File

		f, err := createFile(args.Main)
		if err != nil {
			return err
		}

		files = append(files, f)

		for _, resc := range args.Resources {
			f, err := createFile(resc)
			if err != nil {
				return err
			}

			files = append(files, f)
		}

		cmd := exec.Command("node", f.Name())
		var o bytes.Buffer
		var e bytes.Buffer
		cmd.Stdout = &o
		cmd.Stderr = &e
		err = cmd.Run()

		*reply = RunnerReply{
			OK:        err == nil,
			Output:    o.String(),
			ErrOutput: e.String(),
		}

		cleanFiles(files)
	default:
		return errors.New("We do not support that programming language!")
	}

	return nil
}

func createFile(r File) (*os.File, error) {
	f, err := os.Create(os.TempDir() + r.Name + ".js")
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(r.Contents)
	if err != nil {
		return nil, err
	}

	return f, err
}

func cleanFiles(files []*os.File) {
	for _, f := range files {
		os.Remove(f.Name())
	}
}
