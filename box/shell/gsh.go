package shell

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zooyer/regis/agent/cmd/command"
)

// echo "aaa">>bbb.txt> &1>&2
type execUint struct {
	Name      string
	Args      []string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	Next      *execUint
	Logical   string // && ||
	Append    bool
	Duplicate bool
	//DupFile
}

func doExec(unit execUint) {
	switch unit.Name {
	// 重定向
	case ">>":
		//...

	}

	// 内部命令

	// 外部命令
	exec.Command(unit.Name)
}

func init() {
	// execUnit....
	var units []int

	var wg sync.WaitGroup
	wg.Add(len(units))

	for range units {
		go func() {
			defer wg.Done()

			//
		}()
	}
}

var dupUnit string

func doShell(attr command.Attr) {
	var (
		err   error
		input string
		stdin = bufio.NewReader(attr.Stdin)
	)

	for {
		fmt.Fprint(attr.Stdout, "$ ")

		if input, err = stdin.ReadString('\n'); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			writeError(attr, err)
			continue
		}

		// 判断是否退出
		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// 模拟执行命令（这里只打印用户输入的命令）
		fmt.Printf("Executing: %q\n", input)
	}
}

func Gsh(attr command.Attr, args ...string) (errno int) {
	var (
		err error
		opt Option
		set = flag.NewFlagSet("shell", flag.ContinueOnError)
	)

	// 注册命令行参数
	if err = bindOption(set, &opt); err != nil {
		writeError(attr, err)
		return 1
	}

	// 解析命令行参数
	if err = set.Parse(args); err != nil {
		writeError(attr, err)
		return 2
	}

	if opt.ShellFile != "" {
		if filepath.Base(opt.ShellFile) == opt.ShellFile {
			attr.Stdin = strings.NewReader(opt.ShellFile)
		} else {
			// TODO 判断文件是否有执行权限

			var file *os.File

			// 打开文件
			if file, err = os.Open(opt.ShellFile); err != nil {
				writeError(attr, err)
				return 3
			}

			// 关闭文件
			defer func() {
				if err = file.Close(); err != nil {
					writeError(attr, err)
					errno = 4
				}
			}()

			attr.Stdin = file
		}
	}

	doShell(attr)

	return
}
