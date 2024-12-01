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
		errno  int
		newDir = "../"
		option = types.Option{
			Dir:    "",
			Env:    nil,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	)

	if errno = Cd(option, []string{"cd", newDir}); errno != 0 {
		t.Fatal("Cd failed:", errno)
	}

	if newDir, err = os.Getwd(); err != nil {
		t.Fatal(err)
	}

	if oldDir == newDir {
		t.Fatal("Cd failed, old:", oldDir, "new:", newDir)
	}

	t.Log("old:", oldDir)
	t.Log("new:", newDir)

	if errno = Cd(option, []string{"cd", "-"}); errno != 0 {
		t.Fatal("Cd failed:", errno)
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
		errno  int
		option = types.Option{
			Dir:    "",
			Env:    nil,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
	)

	if errno = Exit(option, []string{"exit", "0"}); errno != 0 {
		t.Fatal("Exit failed:", errno)
	}

	if errno = Exit(option, []string{"exit", "1"}); errno != 1 {
		t.Fatal("Exit failed:", errno)
	}

	if errno = Exit(option, []string{"exit", "99"}); errno != 99 {
		t.Fatal("Exit failed:", errno)
	}
}
