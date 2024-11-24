package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// TokenType Token 类型
type TokenType int

const (
	TokenWord           TokenType = iota // 普通单词（命令或参数）
	TokenPipe                            // 管道符 `|`
	TokenOr                              // 逻辑或 `||`
	TokenAnd                             // 逻辑与 `&&`
	TokenRedirectIn                      // 输入重定向 `<`
	TokenRedirectOut                     // 输出重定向 `>`
	TokenRedirectAppend                  // 追加输出 `>>`
	TokenHereDoc                         // Here Document `<<`
	TokenBackground                      // 后台执行符 `&`
	TokenSemicolon                       // 分号 `;`
	TokenVar                             // 变量 `$XX`
	TokenCmd                             // 命令替换 `$(...)`
)

var tokenSymbols = map[TokenType]string{
	TokenPipe:           "|",
	TokenOr:             "||",
	TokenAnd:            "&&",
	TokenRedirectIn:     "<",
	TokenRedirectOut:    ">",
	TokenRedirectAppend: ">>",
	TokenHereDoc:        "<<",
	TokenBackground:     "&",
	TokenSemicolon:      ";",
}

var (
	symbolTokens    = make(map[string]TokenType)
	symbolMaxLength int
)

var escapeSequences = map[byte]byte{
	'n':  '\n',
	't':  '\t',
	'r':  '\r',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
}

func initSymbolTokens() {
	for token, symbol := range tokenSymbols {
		if symbolTokens[symbol] = token; len(symbol) > symbolMaxLength {
			symbolMaxLength = len(symbol)
		}
	}
}

func init() {
	initSymbolTokens()
}

// Token 结构
type Token struct {
	Type  TokenType
	Value string
}

func readByte(r *bufio.Reader, flag string) (c byte, err error) {
	if c, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("unexpected end of input after %s, err: %w", flag, err)
	}

	return
}

func readEscaped(r *bufio.Reader, flag string) (byte, error) {
	c, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("unexpected end of input after %s, err: %w", flag, err)
	}

	// TODO 解析多字符转义
	switch c {
	case 'u':
	case 'U':
	case 'x':
	case '0':
	}

	if escaped, exists := escapeSequences[c]; exists {
		return escaped, nil
	}

	return c, nil
}

type Lexer interface {
	Token() (*Token, error)
}

type lexer struct {
	reader *bufio.Reader
	buffer []byte
	tokens []Token
}

func NewLexer(reader io.Reader) Lexer {
	return &lexer{
		reader: bufio.NewReader(reader),
		buffer: nil,
	}
}

func (l *lexer) Token() (_ *Token, err error) {
	defer func() {
		if err != nil {
			// 获取缓冲区内容
			buf, _ := l.reader.Peek(l.reader.Buffered())
			context := fmt.Sprintf("near '%s'", string(buf))
			err = fmt.Errorf("%w: %s", err, context)
		}
	}()

	if len(l.tokens) > 0 {
		var token = l.tokens[0]
		l.tokens = l.tokens[1:]
		return &token, nil
	}

	var (
		token   Token
		quotes  []byte
		current = new(strings.Builder)
	)

	var c, nc byte
	for {
		if c, err = l.reader.ReadByte(); err != nil {
			if errors.Is(err, io.EOF) && current.Len() > 0 {
				if len(quotes) > 0 {
					// 检查引号是否闭合
					var quote = quotes[len(quotes)-1]
					return nil, fmt.Errorf("syntax error: unclosed quote '%c' at position %d", quote, current.Len())
				}
				return &Token{Type: TokenWord, Value: current.String()}, nil
			}
			return
		}

		// 在引号内
		if len(quotes) > 0 {
			// 栈顶引号
			var top = quotes[len(quotes)-1]

			// 关闭引号
			if c == top {
				// 弹出栈顶引号
				quotes = quotes[:len(quotes)-1]
				continue
			}

			// 双引号内转义
			if c == '\\' && top == '"' {
				if nc, err = readEscaped(l.reader, "\\"); err != nil {
					return
				}
				current.WriteByte(nc)
				continue
			}

			// 引号内正常字符
			current.WriteByte(c)
			continue
		}

		// 引号外转义符
		if c == '\\' {
			if nc, err = readEscaped(l.reader, "\\"); err != nil {
				return
			}
			current.WriteByte(nc)
			continue
		}

		// 在引号外
		switch c {
		case '\'', '"':
			quotes = append(quotes, c)
		case '\r', '\n':
			if current.Len() > 0 {
				l.tokens = append(l.tokens, Token{Type: TokenSemicolon, Value: tokenSymbols[TokenSemicolon]})
			}
			return &Token{Type: TokenWord, Value: current.String()}, nil
		case ' ', '\t':
			if current.Len() > 0 {
				return &Token{Type: TokenWord, Value: current.String()}, nil
			}
		default:
			for sl := symbolMaxLength - 1; sl >= 0; sl-- {
				var (
					peek   []byte
					symbol = string(c)
				)

				if sl > 0 {
					if peek, err = l.reader.Peek(sl); err != nil {
						if errors.Is(err, io.EOF) {
							continue
						}
						return
					}
					if len(peek) < sl {
						continue
					}
					symbol += string(peek)
				}

				if tokenType, exists := symbolTokens[symbol]; exists {
					if _, err = l.reader.Discard(sl); err != nil {
						return
					}

					token = Token{Type: tokenType, Value: symbol}
					if current.Len() > 0 {
						l.tokens = append(l.tokens, token)
						return &Token{Type: TokenWord, Value: current.String()}, nil
					}

					return &token, nil
				}
			}

			current.WriteByte(c)
		}
	}
}

func ParseTokens(input string) (tokens []Token, err error) {
	var (
		reader = strings.NewReader(input)
		lexer  = NewLexer(reader)
		token  *Token
	)

	for {
		if token, err = lexer.Token(); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			return nil, err
		}

		if token == nil {
			return
		}

		tokens = append(tokens, *token)
	}

	return
}

//type Parser interface {
//	Parse() (ASTNode, error) // 解析生成抽象语法树
//}

//type Executor interface {
//	Execute(node ASTNode) error // 执行语法树
//}
