package shell

import (
	"os"
	"testing"

	"github.com/zooyer/gobox/types"
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
