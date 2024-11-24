package shell

import (
	"os"
	"testing"

	"github.com/zooyer/regis/agent/cmd/command"
)

func TestSh(t *testing.T) {
	var attr = command.Attr{
		Dir:    "",
		Env:    nil,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	Sh(attr, "--debug")
}

func TestShWithEnv(t *testing.T) {
	var arr []string
	arr = arr[1:]
}
