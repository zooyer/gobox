package main

import (
	"os"

	"github.com/zooyer/regis/agent/cmd/command"
	"github.com/zooyer/regis/agent/cmd/command/bool"
)

func main() {
	var attr = command.Attr{
		Dir:    "",
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	os.Exit(bool.True(attr, os.Args[1:]...))
}
