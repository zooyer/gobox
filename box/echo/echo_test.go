package echo

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/zooyer/gobox/box"
)

// systemEchoOutput 使用系统 echo 命令生成期望的输出
func systemEchoOutput(args []string) string {
	cmd := exec.Command("echo", args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func TestEcho(t *testing.T) {
	IsDarwin = runtime.GOOS == "darwin"

	var tests = []struct {
		Name string   // 测试用例名称
		Args []string // Echo 的输入参数
	}{
		{"SimpleString", []string{"hello"}},
		{"MultipleStrings", []string{"hello", "world"}},
		{"Newline", []string{"-e", "hello\\nworld"}},
		{"Tab", []string{"-e", "hello\\tworld"}},
		{"Backslash", []string{"-e", "C:\\\\Windows"}},
		{"EscapeIgnored", []string{"hello\\nworld"}}, // -e 未启用
		{"HelpOption", []string{"-h"}},
		{"VersionOption", []string{"-v"}},
		{"EmptyString", []string{""}},
		{"NoArgs", []string{}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// 使用系统 echo 获取期望结果
			expected := systemEchoOutput(test.Args)

			// 调用自定义 Echo 函数
			var (
				stdout, stderr bytes.Buffer
				option         = box.Option{
					Args:   append([]string{"echo"}, test.Args...),
					Stdout: &stdout,
					Stderr: &stderr,
				}
			)

			// 执行测试
			if errno := Echo(context.Background(), option); errno != 0 {
				t.Fatalf("[%s] Unexpected errno: got %d, want 0", test.Name, errno)
			}

			// 比较标准输出
			if strings.TrimSpace(stdout.String()) != strings.TrimSpace(expected) {
				t.Errorf("[%s] Unexpected stdout:\nGot:\n%s\nWant:\n%s", test.Name, stdout.String(), expected)
			}

			// 标准错误输出应为空
			if stderr.Len() > 0 {
				t.Errorf("[%s] Unexpected stderr: got %s, want empty", test.Name, stderr.String())
			}
		})
	}
}
