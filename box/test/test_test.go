package test

import (
	"github.com/zooyer/gobox/types"
	"testing"
)

func TestNew(t *testing.T) {
	var option = types.Option{
		Dir:     "",
		Args:    nil,
		Envs:    nil,
		Stdin:   nil,
		Stdout:  nil,
		Stderr:  nil,
		Syscall: types.Syscall{},
	}

	var test = New(option)
	test.Main(nil)
	test.Kill()
	test.Wait()
}
