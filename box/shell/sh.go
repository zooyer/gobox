package shell

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/zooyer/regis/agent/cmd/command"
	"io"
	"strings"
)

const version = "1.0.0"

const usage = `%s, version %s
Usage:	%s [GNU long option] [option] ...
	%s [GNU long option] [option] script-file ...

GNU long options:
	--debug
	--help
	--version
	--login
	-i (interactive mode)
Shell options:
	Type '%s -c "help set"' for more information about shell options.
`

/*
GNU long options:


	--rcfile
	--restricted
	--verbose
	--version
Shell options:
	-ilrsD or -c command or -O shopt_option		(invocation only)
	-abefhkmnptuvxBCEHPT or -o option
*/

type ShRunOption struct {
	Interactive bool `json:"i"`
	Login       bool `json:"l"`
	Restricted  bool `json:"r"`
	Verbose     bool `json:"v"`
	NoClobber   bool `json:"C"`
	Debug       bool `json:"D"`
	NoEditing   bool `json:"n"`
}

type ShConfigOption struct {
	AllExport   bool   `json:"a"`
	BraceExpand bool   `json:"B"`
	EmacsEdit   bool   `json:"e"`
	NoBuiltin   bool   `json:"b"`
	ShellFile   string `json:"c"`
	NoProfile   bool   `json:"P"`
}

type ShOption struct {
	GNUOption
	ShRunOption
	ShConfigOption
}

func help(name string, attr command.Attr) {
	_, _ = fmt.Fprintf(attr.Stdout, usage, name, version, name, name, name)
}

func Sh(attr command.Attr, args ...string) (errno int) {
	var (
		err error
		opt ShOption
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

	data, err := json.MarshalIndent(opt, "", "  ")
	if err != nil {
		writeError(attr, err)
		return 2
	}

	fmt.Println(string(data))

	// TODO 判断option
	switch {
	case opt.Help:
		help("sh", attr)
		return 0
	}

	var (
		line  string
		stdin      = bufio.NewReader(attr.Stdin)
		delim byte = '\n'
		//end    string
		buffer bytes.Buffer
	)
	for {
		if buffer.Len() > 0 {
			fmt.Fprint(attr.Stdout, "> ")
		} else {
			fmt.Fprint(attr.Stdout, "$ ")
		}
		//stdin.ReadLine()
		if line, err = stdin.ReadString(delim); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			writeError(attr, err)

			continue
		}

		buffer.WriteString(line)

		//if needsContinuation(line) {
		//	continue
		//}

		// 执行命令
		var input = buffer.String()
		input = strings.TrimSpace(input)
		buffer.Reset()

		// 判断是否退出
		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// 模拟执行命令（这里只打印用户输入的命令）
		fmt.Printf("Executing: %q\n", input)
	}

	return
}