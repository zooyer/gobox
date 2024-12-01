package shell

import (
	"bufio"
	"context"
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
	TokenHeredoc                         // Here Document `<<`
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
	TokenHeredoc:        "<<",
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

func readBytes(r *bufio.Reader, delim byte, flag string) (bytes []byte, err error) {
	if bytes, err = r.ReadBytes(delim); err != nil {
		return nil, fmt.Errorf("unexpected end of input after %s, err: %w", flag, err)
	}

	return
}

func readEscaped(r *bufio.Reader, flag string) (c byte, err error) {
	if c, err = readByte(r, flag); err != nil {
		return
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

func readDelimiter(r *bufio.Reader, delim string, ident bool, flag string) (_ string, err error) {
	var (
		sb     strings.Builder
		data   []byte
		suffix = "\n" + delim + "\n"
	)

	for {
		if data, err = readBytes(r, '\n', flag); err != nil {
			return
		}

		if !ident {
			sb.Write(data)
		} else {
			var front = true
			for _, b := range data {
				switch b {
				case '\n':
					front = true
				case '\t':
					if front {
						continue
					}
				default:
					front = false
				}

				sb.WriteByte(b)
			}
		}

		// TODO 考虑\r,\r\n，ident去除每行前面\t
		if strings.HasSuffix(sb.String(), suffix) {
			return strings.TrimSuffix(sb.String(), suffix), nil
		}
	}
}

type Lexer struct {
	err    error
	reader *bufio.Reader
	tokens chan Token
}

func (l *Lexer) getValue(sb *strings.Builder) string {
	if sb == nil {
		return ""
	}

	defer sb.Reset()

	return sb.String()
}

func (l *Lexer) inputToken(tokenType TokenType, tokenValue string) {
	l.tokens <- Token{Type: tokenType, Value: tokenValue}
}

func (l *Lexer) inputWordToken(sb *strings.Builder) {
	if sb != nil && sb.Len() > 0 {
		l.inputToken(TokenWord, sb.String())
		sb.Reset()
	}
}

func (l *Lexer) isRun(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return true
	}
}

func (l *Lexer) Token() <-chan Token {
	return l.tokens
}

func (l *Lexer) Run(ctx context.Context) (err error) {
	var (
		word    = new(strings.Builder)
		quotes  []byte
		ident   bool
		delim   string
		heredoc bool
	)

	if l.err != nil {
		return l.err
	}

	defer func() {
		close(l.tokens)
		if err != nil {
			l.err = err
		}
	}()

	var c, nc byte
	for {
		if !l.isRun(ctx) {
			return
		}

		if c, err = l.reader.ReadByte(); err != nil {
			if errors.Is(err, io.EOF) {
				break
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
				word.WriteByte(nc)
				continue
			}

			// 引号内正常字符
			word.WriteByte(c)
			continue
		}

		// 引号外转义符
		if c == '\\' {
			if nc, err = readEscaped(l.reader, "\\"); err != nil {
				return
			}
			word.WriteByte(nc)
			continue
		}

		// 在引号外
		switch c {
		case '\'', '"':
			quotes = append(quotes, c)
		case ' ', '\t', '\r', '\n':
			if heredoc && delim == "" {
				delim = l.getValue(word)
			} else {
				l.inputWordToken(word)
			}

			if c == '\r' || c == '\n' {
				if heredoc {
					if delim == "" {
						return fmt.Errorf("syntax error: heredoc delim is null, position %d", word.Len())
					}

					var value string
					if value, err = readDelimiter(l.reader, delim, ident, "<<"); err != nil {
						return
					}

					heredoc = false
					delim = ""
					ident = false

					// TODO 考虑合并成一个还是多个
					l.inputToken(TokenHeredoc, value)
					//l.inputToken(TokenHeredoc, "<<") // << or <<-
					//l.inputWordToken(word)
				}
				l.inputToken(TokenSemicolon, tokenSymbols[TokenSemicolon])
			}
		case '-':
			if heredoc && word.Len() == 0 && delim == "" {
				ident = true
				continue
			}
			fallthrough
		default:
			var hasSymbol bool
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

					hasSymbol = true
					l.inputWordToken(word)

					if tokenType == TokenHeredoc {
						heredoc = true
					} else {
						l.inputToken(tokenType, symbol)
					}

					break
				}
			}

			if !hasSymbol {
				word.WriteByte(c)
			}
		}
	}

	// 检查引号是否闭合
	if word.Len() > 0 && len(quotes) > 0 {
		var quote = quotes[len(quotes)-1]
		return fmt.Errorf("syntax error: unclosed quote '%c' at position %d", quote, word.Len())
	}

	l.inputWordToken(word)

	// 检查heredoc完整结束
	if heredoc {
		return fmt.Errorf("unexpected end of input after heredoc delim")
	}

	return nil
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		tokens: make(chan Token, 1024),
	}
}

func ParseTokens(input string) (tokens []Token, err error) {
	var (
		reader = strings.NewReader(input)
		lexer  = NewLexer(reader)
		token  = lexer.Token()
	)

	go func() {
		if err = lexer.Run(context.Background()); err != nil {
			return
		}
	}()

	for tk := range token {
		tokens = append(tokens, tk)
	}

	return
}
