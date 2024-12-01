package shell

import (
	"os"
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestCd(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	var (
		code   int
		newDir = "../"
		option = types.Option{
			Dir:    "",
			Env:    nil,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	)

	if code = Cd(option, []string{"cd", newDir}); code != 0 {
		t.Fatal("Cd failed:", code)
	}

	if newDir, err = os.Getwd(); err != nil {
		t.Fatal(err)
	}

	if oldDir == newDir {
		t.Fatal("Cd failed, old:", oldDir, "new:", newDir)
	}

	t.Log("old:", oldDir)
	t.Log("new:", newDir)

	if code = Cd(option, []string{"cd", "-"}); code != 0 {
		t.Fatal("Cd failed:", code)
	}

	if newDir, err = os.Getwd(); err != nil {
		t.Fatal(err)
	}

	if oldDir != newDir {
		t.Fatal("Cd failed, old:", oldDir, "new:", newDir)
	}

	t.Log("old:", oldDir)
	t.Log("new:", newDir)
}

func TestExit(t *testing.T) {
	var (
		code   int
		option = types.Option{
			Dir:    "",
			Env:    nil,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	)

	if code = Exit(option, []string{"exit", "0"}); code != 0 {
		t.Fatal("Exit failed:", code)
	}

	if code = Exit(option, []string{"exit", "1"}); code != 1 {
		t.Fatal("Exit failed:", code)
	}

	if code = Exit(option, []string{"exit", "99"}); code != 99 {
		t.Fatal("Exit failed:", code)
	}
}
