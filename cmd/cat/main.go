package main

import (
	"os"

	"github.com/zooyer/gobox/box/cat"
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

	os.Exit(cat.New(opt).Main(os.Args))
}
