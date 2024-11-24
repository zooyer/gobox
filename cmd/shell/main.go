package main

import (
	"os"

	"github.com/zooyer/regis/agent/cmd/command"
	"github.com/zooyer/regis/agent/cmd/command/shell"
)

func main() {
	var attr = command.Attr{
		Dir:    "",
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	os.Exit(shell.Sh(attr, os.Args[1:]...))
}
