package types

import (
	"io"
	"os"
)

type Option struct {
	Dir string
	Env []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Process interface {
	Main(args []string) (errno int)
	Kill()
	Wait()
	Signal(signal os.Signal)
}

type NewFunc func(option Option) Process

type MainFunc func(opt Option, args []string) (errno int)

// head/tail/wc/nl/tac/yes/tee/sort
/*
basename：提取文件名。
dirname：提取文件路径。
pwd：打印当前工作目录。
*/
