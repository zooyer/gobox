package main

import (
	"context"
	"os"

	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/box/echo"
)

func main() {
	var opt = box.Option{
		Dir:    "",
		Envs:   os.Environ(),
		Args:   os.Args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	os.Exit(echo.Echo(context.Background(), opt))
}
