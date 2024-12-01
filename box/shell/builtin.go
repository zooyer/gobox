package shell

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/zooyer/gobox/types"
)

func expandHome(path string) (_ string, err error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	var (
		dirs = filepath.SplitList(path)
		dir0 = dirs[0]
		usr  *user.User
	)

	if len(dir0) > 1 {
		usr, err = user.Lookup(dir0[1:])
	} else {
		usr, err = user.Current()
	}
	if err != nil {
		return
	}

	dirs[0] = usr.HomeDir

	var dir = filepath.Join(dirs...)

	return dir, nil
}

func expandLast(opt types.Option, path string) (_ string, err error) {
	if !strings.HasPrefix(path, "-") {
		return path, nil
	}

	var index int
	if len(path) > 1 {
		if index, err = strconv.Atoi(path[1:]); err != nil {
			return
		}
	}

	// TODO 获取历史工作目录
	_ = index

	var dir = os.Getenv("OLDPWD")
	if dir == "" {
		return ".", nil
	}

	return dir, nil
}

func compareDir(old, new string) (same bool, err error) {
	if old, err = filepath.Abs(old); err != nil {
		return
	}

	if new, err = filepath.Abs(new); err != nil {
		return
	}

	return old == new, nil
}

func Cd(opt types.Option, args []string) (errno int) {
	var (
		err      error
		end      bool
		set      bool
		dir      string
		escape   = true
		logical  bool
		physical bool
	)

	for _, arg := range args[1:] {
		if end {
			if set {
				writeError(opt, fmt.Errorf("too many arguments"))
				return 1
			}
			dir = arg
			set = true
			escape = false
			continue
		}

		switch arg {
		case "--":
			end = true
		case "-P":
			physical = true
		case "-L":
			logical = true
		default:
			if set {
				writeError(opt, fmt.Errorf("too many arguments"))
				return 1
			}
			dir = arg
			set = true
			escape = true
		}
	}

	if dir == "" {
		dir = "~"
	}

	switch dir[0] {
	case '-':
		if escape {
			if dir, err = expandLast(opt, dir); err != nil {
				writeError(opt, err)
				return 2
			}
		}
	case '~':
		if dir, err = expandHome(dir); err != nil {
			writeError(opt, err)
			return 3
		}
	}

	_ = logical
	if physical {
		dir, err = filepath.EvalSymlinks(dir)
		if err != nil {
			writeError(opt, err)
			return 4
		}
	} else {
		dir = filepath.Clean(dir)
	}

	var old string
	if old, err = os.Getwd(); err != nil {
		writeError(opt, err)
		return 5
	}

	var same bool
	if same, err = compareDir(old, dir); err != nil {
		writeError(opt, err)
		return 6
	}
	if same {
		return
	}

	if err = os.Chdir(dir); err != nil {
		writeError(opt, err)
		return 7
	}

	// TODO 保存到历史工作目录，相同目录则不保存
	if err = os.Setenv("OLDPWD", old); err != nil {
		writeError(opt, err)
		return 8
	}

	if err = os.Setenv("PWD", dir); err != nil {
		writeError(opt, err)
		return 9
	}

	return 0
}

const exitUsage = `exit: exit [n]
    Exit the shell.
    
    Exits the shell with a status of N.  If N is omitted, the exit status
    is that of the last command executed.
`

func Exit(opt types.Option, args []string) (errno int) {
	// TODO 获取上个命令的退出码

	var err error

	for _, arg := range args[1:] {
		switch arg {
		case "-h", "--help":
			_, _ = fmt.Fprint(opt.Stdout, exitUsage)

			return
		}

		if errno, err = strconv.Atoi(arg); err != nil {
			writeError(opt, err)
			return -1
		}
	}

	if testing.Testing() {
		syscall.Exit(errno)
	}

	os.Exit(errno)

	return
}
