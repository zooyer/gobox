package echo

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/zooyer/gobox/box"
)

// 转义字符映射
var escapeMap = map[byte]string{
	'n':  "\n",
	't':  "\t",
	'r':  "\r",
	'v':  "\v",
	'b':  "\b",
	'f':  "\f",
	'\\': "\\",
	'\'': "'",
	'"':  "\"",
}

// interpretEscapes 使用 map 解析转义字符
func interpretEscapes(input string) string {
	var builder strings.Builder
	length := len(input)

	for i := 0; i < length; i++ {
		if input[i] == '\\' && i+1 < length { // 检测到转义符
			if replacement, exists := escapeMap[input[i+1]]; exists {
				builder.WriteString(replacement)
				i++ // 跳过转义字符
			} else {
				// 未识别的转义序列，保留原样
				builder.WriteByte('\\')
				builder.WriteByte(input[i+1])
				i++
			}
		} else {
			builder.WriteByte(input[i]) // 普通字符直接追加
		}
	}

	return builder.String()
}

func writeError(opt box.Option, err error) {
	_, _ = fmt.Fprintln(opt.Stderr, "echo:", err)
}

var IsDarwin = runtime.GOOS == "darwin"

func Echo(ctx context.Context, opt box.Option) (errno int) {
	var (
		e, n bool
		end  bool
		out  = make([]string, 0, len(opt.Args))
	)

	for _, arg := range opt.Args[1:] {
		if end {
			out = append(out, arg)
			continue
		}

		switch arg {
		case "-n":
			n = true
		case "--":
			end = true
		case "-e":
			if !IsDarwin {
				e = true
				break
			}
			fallthrough
		default:
			end = true
			out = append(out, arg)
		}
	}

	var (
		result  = strings.Join(out, " ")
		doPrint = fmt.Fprintln
	)

	if e {
		result = interpretEscapes(result)
	}

	if n {
		doPrint = fmt.Fprint
	}

	if _, err := doPrint(opt.Stdout, result); err != nil {
		errno = 1
		writeError(opt, err)
	}

	return
}