package shell

import (
	"github.com/zooyer/gobox/types"
	"os"
	"testing"
)

func TestSh(t *testing.T) {
	var opt = types.Option{
		Dir:    "",
		Env:    nil,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	Sh(opt, "--debug")
}

func TestShWithEnv(t *testing.T) {
	var arr []string
	arr = arr[1:]
}
