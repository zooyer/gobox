package main

import (
	"github.com/zooyer/regis/agent/cmd/command/cat"
	"os"

	"github.com/zooyer/regis/agent/cmd/command"
)

func main() {
	var attr = command.Attr{
		Dir:    "",
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	os.Exit(cat.Cat(attr, os.Args[1:]...))
}
