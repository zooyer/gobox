package box

import (
	"context"
	"io"
)

type Option struct {
	Dir    string
	Args   []string
	Envs   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Main func(ctx context.Context, opt Option) (errno int)

// head/tail/wc/nl/tac/yes/tee/sort
/*
basename：提取文件名。
dirname：提取文件路径。
pwd：打印当前工作目录。
*/
