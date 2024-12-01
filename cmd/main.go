package main

import (
	"fmt"
	"os"

	"github.com/zooyer/gobox/box/cmd"
	"github.com/zooyer/gobox/types"
)

func main() {
	var opt = types.Option{
		Dir:    "",
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	var fn = cmd.New(os.Args[1])
	if fn == nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s is not a command. See '%s --help'.", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	os.Exit(fn(opt).Main(os.Args[1:]))
}
