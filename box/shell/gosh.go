package shell

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/box/cmd"
	"github.com/zooyer/gobox/types"
)

type Gosh struct {
	box.Process
	Builtin map[string]types.MainFunc // 内置命令
	Command map[string]types.NewFunc  // 系统命令
}

func (sh *Gosh) ps1(option types.Option) {
	var (
		err  error
		dir  string
		usr  *user.User
		host string
	)

	if dir, err = os.Getwd(); err != nil {
		dir = ""
	}

	if usr, err = user.Current(); err != nil {
		usr = &user.User{
			Uid:      "0",
			Gid:      "0",
			Username: "gosh",
			Name:     "GOSH",
			HomeDir:  dir,
		}
	}

	if dir == usr.HomeDir {
		dir = "~"
	}

	dir = filepath.Base(dir)

	if host, err = os.Hostname(); err != nil {
		host = "gosh"
	}

	_, _ = fmt.Fprintf(option.Stdout, "%s@%s:%s$ ", host, usr.Username, dir)
}

func (sh *Gosh) Run(stdin io.Reader, option types.Option) (code int, err error) {
	var (
		wg      sync.WaitGroup
		ctx     = context.Background()
		lexer   = NewLexer(stdin)
		parser  = NewParser(lexer.Token())
		command = parser.Command()
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		if err = lexer.Run(ctx); err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()

		if err = parser.Run(ctx); err != nil {
			return
		}
	}()

	for command := range command {
		if code, err = sh.Exec(&command, option); err != nil {
			writeError(option, err)
			return 2, err
		}
	}

	return
}

func (sh *Gosh) Exec(command *Command, option types.Option) (code int, err error) {
	if command == nil {
		return 0, errors.New("nil command")
	}

	var (
		wg         sync.WaitGroup
		stdin      *os.File
		stdout     *os.File
		thisOption = option
	)

	// 执行 管道后命令（并行执行）
	if command.Pipe != nil {
		var pipeOption = option

		// 普通管道
		//io.Pipe()

		// 进程间管道
		if stdin, stdout, err = os.Pipe(); err != nil {
			return
		}

		// 管道重定向
		pipeOption.Stdin, thisOption.Stdout = stdin, stdout

		var (
			pipeErr  error
			pipeCode int
		)

		// 使用管道的退出状态码
		defer func() {
			if code = pipeCode; err == nil && pipeErr != nil {
				err = pipeErr
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer deferClose(&err, stdin.Close)
			if pipeCode, pipeErr = sh.Exec(command.Pipe, pipeOption); err != nil {
				return
			}
		}()
	}

	// 重定向 文件输入
	if command.Input != "" {
		var input *os.File
		if input, err = os.Open(command.Input); err != nil {
			return
		}
		thisOption.Stdin = input
		defer deferClose(&err, input.Close)
	}

	// 重定向 文档输入
	if command.Heredoc != "" {
		thisOption.Stdin = strings.NewReader(command.Heredoc)
	}

	// 重定向 文件输出
	if command.Output != "" {
		var output *os.File
		if output, err = os.OpenFile(command.Output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			return
		}
		thisOption.Stdout = output
		defer deferClose(&err, output.Close)
	}

	// 重定向 文件追加
	if command.Append != "" {
		var output *os.File
		if output, err = os.OpenFile(command.Append, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
			return
		}
		thisOption.Stdout = output
		defer deferClose(&err, output.Close)
	}

	// 执行 当前命令
	switch {
	case sh.Builtin != nil && sh.Builtin[command.Path] != nil:
		var cmd = sh.Builtin[command.Path]

		if command.Background {
			go cmd(thisOption, command.CmdArgs())
		}

		code = cmd(thisOption, command.CmdArgs())
	case sh.Command != nil && sh.Command[command.Path] != nil:
		var cmd = sh.Command[command.Path](thisOption)

		if command.Background {
			go cmd.Main(command.CmdArgs())
		}

		code = cmd.Main(command.CmdArgs())
	default:
		var cmd = exec.Command(command.Path, command.Args...)
		cmd.Dir = thisOption.Dir
		cmd.Env = thisOption.Env
		cmd.Stdin = thisOption.Stdin
		cmd.Stdout = thisOption.Stdout
		cmd.Stderr = thisOption.Stderr

		if err = cmd.Start(); err != nil {
			return
		}

		if command.Background {
			go func() { _ = cmd.Wait() }()
		} else {
			if err = cmd.Wait(); err != nil {
				return
			}

			code = cmd.ProcessState.ExitCode()
		}
	}

	if stdout != nil {
		deferClose(&err, stdout.Close)
	}

	// 执行 与
	if command.And != nil && code == 0 {
		if code, err = sh.Exec(command.And, option); err != nil {
			return
		}
	}

	// 执行 或
	if command.Or != nil && code != 0 {
		if code, err = sh.Exec(command.Or, option); err != nil {
			return
		}
	}

	wg.Wait()

	return
}

func (sh *Gosh) Main(args []string) (code int) {
	var (
		err    error
		opt    Option
		set    = flag.NewFlagSet("gosh", flag.ContinueOnError)
		option = sh.Option
	)

	// 注册命令行参数
	if err = bindOption(set, &opt); err != nil {
		writeError(option, err)
		return 1
	}

	// 解析命令行参数
	if err = set.Parse(args); err != nil {
		writeError(option, err)
		return 2
	}

	if opt.ShellFile != "" {
		if filepath.Base(opt.ShellFile) == opt.ShellFile {
			option.Stdin = strings.NewReader(opt.ShellFile)
		} else {
			// TODO 判断文件是否有执行权限

			var file *os.File

			// 打开文件
			if file, err = os.Open(opt.ShellFile); err != nil {
				writeError(option, err)
				return 3
			}

			// 关闭文件
			defer func() {
				if err = file.Close(); err != nil {
					writeError(option, err)
					code = 4
				}
			}()

			option.Stdin = file
		}
	}

	var (
		ctx     = context.Background()
		errs    = make(chan error, 2)
		lexer   = NewLexer(option.Stdin)
		parser  = NewParser(lexer.Token())
		command = parser.Command()
	)

	go func() { errs <- lexer.Run(ctx) }()
	go func() { errs <- parser.Run(ctx) }()

	sh.ps1(option)
	for command := range command {
		if code, err = sh.Exec(&command, option); err != nil {
			writeError(option, err)
			return 2
		}
		sh.ps1(option)
	}

	for err = range errs {
		if err != nil {
			writeError(option, err)
			return 3
		}
	}

	return
}

func NewGosh(opt types.Option) *Gosh {
	var (
		command = cmd.Cmd()
		builtin = map[string]types.MainFunc{
			"cd":   Cd,
			"exit": Exit,
		}
	)

	command["gosh"] = func(option types.Option) types.Process {
		return NewGosh(option)
	}

	return &Gosh{
		Process: box.Process{
			Option: opt,
		},
		Builtin: builtin,
		Command: command,
	}
}
