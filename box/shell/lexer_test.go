package shell

import (
	"reflect"
	"testing"
)

func TestParseTokens(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Token
		err      bool
	}{
		{
			// 普通命令
			input: "echo hello",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
			},
			err: false,
		},
		{
			// 带引号的命令
			input: "echo \"hello world\"",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello world"},
			},
			err: false,
		},
		{
			// 单引号
			input: "echo 'hello world'",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello world"},
			},
			err: false,
		},
		{
			// 带转义字符的引号
			input: "echo \"hello\\\"world\"",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello\"world"},
			},
			err: false,
		},
		{
			// 管道符
			input: "echo hello | grep world",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenPipe, Value: "|"},
				{Type: TokenWord, Value: "grep"},
				{Type: TokenWord, Value: "world"},
			},
			err: false,
		},
		{
			// 逻辑运算符 AND
			input: "echo hello && echo world",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenAnd, Value: "&&"},
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "world"},
			},
			err: false,
		},
		{
			// 逻辑运算符 OR
			input: "echo hello || echo world",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenOr, Value: "||"},
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "world"},
			},
			err: false,
		},
		{
			// 重定向
			input: "echo hello > output.txt",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenRedirectOut, Value: ">"},
				{Type: TokenWord, Value: "output.txt"},
			},
			err: false,
		},
		{
			// 追加输出重定向
			input: "echo hello >> output.txt",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenRedirectAppend, Value: ">>"},
				{Type: TokenWord, Value: "output.txt"},
			},
			err: false,
		},
		{
			// 输入重定向
			input: "cat < input.txt",
			expected: []Token{
				{Type: TokenWord, Value: "cat"},
				{Type: TokenRedirectIn, Value: "<"},
				{Type: TokenWord, Value: "input.txt"},
			},
			err: false,
		},
		{
			// here doc 重定向
			input: "cat << EOF\nHello World\nEOF\n",
			expected: []Token{
				{Type: TokenWord, Value: "cat"},
				{Type: TokenHeredoc, Value: "Hello World"},
				{Type: TokenSemicolon, Value: ";"},
			},
			err: false,
		},
		{
			// 后台执行符
			input: "echo hello &",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenBackground, Value: "&"},
			},
			err: false,
		},
		{
			// 分号
			input: "echo hello; echo world",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenSemicolon, Value: ";"},
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "world"},
			},
			err: false,
		},
		{
			// 错误的引号（未闭合）
			input: "echo \"hello",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
			},
			err: true,
		},
		{
			// 错误的多重分号
			input: "echo hello ;; echo world",
			expected: []Token{
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenSemicolon, Value: ";"},
				{Type: TokenSemicolon, Value: ";"},
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "world"},
			},
			err: false,
		},
		{
			input: "cat << EOF &\nBackground task\nEOF\necho \"HereDoc submitted\"",
			expected: []Token{
				{Type: TokenWord, Value: "cat"},
				{Type: TokenBackground, Value: "&"},
				{Type: TokenHeredoc, Value: "Background task"},
				{Type: TokenSemicolon, Value: ";"},
				{Type: TokenWord, Value: "echo"},
				{Type: TokenWord, Value: "HereDoc submitted"},
			},
			err: false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens, err := ParseTokens(test.input)
			if (err != nil) != test.err {
				t.Errorf("expected error: %v, got: %v", test.err, err)
			}
			if !reflect.DeepEqual(tokens, test.expected) {
				t.Errorf("expected: %v, got: %v", test.expected, tokens)
			}
		})
	}
}
