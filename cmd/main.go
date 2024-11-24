package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zooyer/gobox"
	"github.com/zooyer/gobox/box"
)

func main() {
	var opt = box.Option{
		Dir:    "",
		Envs:   os.Environ(),
		Args:   os.Args[1:],
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	var fn = gobox.Get(os.Args[1])
	if fn == nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s is not a command. See '%s --help'.", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	os.Exit(fn(context.Background(), opt))
}
