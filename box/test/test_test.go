package test

import (
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestNew(t *testing.T) {
	var option = types.Option{
		Dir:    "",
		Env:    nil,
		Stdin:  nil,
		Stdout: nil,
		Stderr: nil,
	}

	var test = New(option)
	test.Main(nil)
	test.Kill()
	test.Wait()
}
