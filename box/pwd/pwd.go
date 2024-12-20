package pwd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/types"
)

const version = "pwd (by golang) 1.0.0"

const usage = `Usage: pwd [OPTION]...
Print the full filename of the current working directory.

  -L, --logical   use PWD from environment, even if it contains symlinks
  -P, --physical  avoid all symlinks
      --help      display this help and exit
      --version   output version information and exit
`

type Pwd struct {
	box.Process
}

type Option struct {
	Logical  bool
	Physical bool
	Help     bool
	Version  bool
}

func writeError(opt types.Option, err error) {
	_, _ = fmt.Fprintln(opt.Stderr, "pwd:", err)
}

// 解析和处理选项
func parse(args []string) (option Option, err error) {
	for _, arg := range args {
		switch arg {
		case "-L", "--logical":
			option.Logical = true
			option.Physical = false
		case "-P", "--physical":
			option.Logical = false
			option.Physical = true
		case "--help":
			option.Help = true
		case "--version":
			option.Version = true
		default:
			err = errors.New("invalid option " + arg)
			return
		}
	}

	// 如果既没有 -L 也没有 -P，默认 -P
	if !option.Logical && !option.Physical {
		option.Physical = true
	}

	return
}

func (p *Pwd) Main(args []string) (code int) {
	option, err := parse(args[1:])
	if err != nil {
		writeError(p.Option, err)
		return 1
	}

	var cwd string

	switch {
	case option.Help:
		_, _ = fmt.Fprint(p.Option.Stdout, usage)
		return
	case option.Version:
		_, _ = fmt.Fprint(p.Option.Stdout, version)
		return
	case option.Physical:
		if cwd, err = os.Getwd(); err == nil {
			cwd, err = filepath.EvalSymlinks(cwd)
		}
	case option.Logical:
		if cwd = os.Getenv("PWD"); cwd == "" {
			cwd, err = os.Getwd()
		}
	default:
		cwd, err = os.Getwd()
	}

	if cwd == "" && err == nil {
		cwd, err = os.Getwd()
	}

	if err != nil {
		writeError(p.Option, err)
		return 2
	}

	_, _ = fmt.Fprintln(p.Option.Stdout, cwd)

	return
}

func New(opt types.Option) types.Process {
	return &Pwd{
		Process: box.Process{
			Option: opt,
		},
	}
}
