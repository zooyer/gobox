package cat

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestCat(t *testing.T) {
	var tests = []struct {
		Name         string // 用例名称
		Args         []string
		Opts         types.Option // 输入选项
		ExpectStdout string       // 期望 stdout 包含的内容（可为空）
		ExpectStderr bool         // 期望 stderr 是否有值
		Code         int          // 期望的错误码
		Setup        func()       // 每个测试用例的预处理
		Cleanup      func()       // 每个测试用例的清理逻辑
	}{
		{
			Name: "Empty file", // 测试空文件
			Args: []string{"empty.txt"},
			Opts: types.Option{Stdin: nil, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
			Setup: func() {
				// 创建一个空文件
				if _, err := os.Create("empty.txt"); err != nil {
					t.Fatal("Setup failed:", err)
				}
			},
			Cleanup: func() {
				_ = os.Remove("empty.txt")
			},
			ExpectStdout: "",
			ExpectStderr: false,
			Code:         0,
		},
		{
			Name:         "Non-existent file", // 测试不存在的文件
			Args:         []string{"nonexistent.txt"},
			Opts:         types.Option{Stdin: nil, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
			Setup:        func() {}, // 不需要额外设置
			Cleanup:      func() {},
			ExpectStdout: "",
			ExpectStderr: true, // 预期 stderr 有值
			Code:         3,
		},
		{
			Name: "Help option", // 测试 -h 选项
			Args: []string{"-h"},
			Opts: types.Option{Stdin: nil, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},

			Setup:        func() {},
			Cleanup:      func() {},
			ExpectStdout: "Usage:", // 只需包含部分内容
			ExpectStderr: false,
			Code:         0,
		},
		{
			Name:         "Version option", // 测试 -v 选项
			Args:         []string{"-v"},
			Opts:         types.Option{Stdin: nil, Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}},
			Setup:        func() {},
			Cleanup:      func() {},
			ExpectStdout: "cat: (by goland)", // 部分内容即可
			ExpectStderr: false,
			Code:         0,
		},
		{
			Name: "Read stdin", // 测试标准输入
			Args: []string{"-"},
			Opts: types.Option{
				Stdin:  bytes.NewBufferString("stdin content"),
				Stdout: &bytes.Buffer{},
				Stderr: &bytes.Buffer{},
			},
			Setup:        func() {},
			Cleanup:      func() {},
			ExpectStdout: "stdin content",
			ExpectStderr: false,
			Code:         0,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Setup != nil {
				test.Setup()
			}
			if test.Cleanup != nil {
				defer test.Cleanup()
			}

			// 执行测试
			stdout := test.Opts.Stdout.(*bytes.Buffer)
			stderr := test.Opts.Stderr.(*bytes.Buffer)

			// 验证错误码
			if code := New(test.Opts).Main(append([]string{"echo"}, test.Args...)); code != test.Code {
				t.Fatalf("Unexpected code: got %d, want %d", code, test.Code)
			}

			// 验证 stdout
			if test.ExpectStdout != "" && !strings.Contains(stdout.String(), test.ExpectStdout) {
				t.Fatalf("Unexpected stdout: got %q, want to contain %q", stdout.String(), test.ExpectStdout)
			}

			// 验证 stderr
			if test.ExpectStderr && stderr.Len() == 0 {
				t.Fatalf("Expected stderr to have content, but it was empty")
			} else if !test.ExpectStderr && stderr.Len() > 0 {
				t.Fatalf("Expected stderr to be empty, but got %q", stderr.String())
			}
		})
	}
}
