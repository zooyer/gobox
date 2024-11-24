package main

import (
	"fmt"
	"strings"
)

func appendWord(tokens []Token, sb *strings.Builder) []Token {
	if sb.Len() == 0 {
		return tokens
	}

	defer sb.Reset()

	return append(tokens, Token{Type: TokenWord, Value: sb.String()})
}

func getLines(index int, input string) (line, col int) {
	if index > len(input) {
		return -1, -1
	}

	line, col = 1, 1

	for i := 0; i < index; i++ {
		switch input[i] {
		case '\n':
			line++
			col = 1
		case '\r':
			continue
		default:
			col++
		}
	}

	return
}

func Tokenize(input string) (tokens []Token, err error) {
	var (
		quotes  []byte
		current = new(strings.Builder)
	)

	for i := 0; i < len(input); i++ {
		var c = input[i]

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
				if i+1 >= len(input) {
					return nil, fmt.Errorf("unexpected end of input after \\")
				}

				i++
				current.WriteRune(rune(input[i]))
				continue
			}

			// 引号内正常字符
			current.WriteByte(c)
			continue
		}

		// 引号外转义符
		if c == '\\' {
			if i+1 >= len(input) {
				return nil, fmt.Errorf("unexpected end of input after \\")
			}
			i++
			current.WriteByte(input[i])
			continue
		}

		// 在引号外
		switch c {
		case '\'', '"':
			quotes = append(quotes, c)
		case ' ', '\t', '\r', '\n':
			tokens = appendWord(tokens, current)
		default:
			var isSymbol bool
			for sl := symbolMaxLength; sl > 0; sl-- {
				if i+sl > len(input) {
					continue
				}

				var symbol = input[i : i+sl]
				if token, exists := symbolTokens[symbol]; exists {
					i += sl - 1
					isSymbol = true
					tokens = appendWord(tokens, current)
					tokens = append(tokens, Token{Type: token, Value: symbol})
					break
				}
			}

			if !isSymbol {
				current.WriteByte(c)
			}
		}
	}

	// 检查引号是否闭合
	if len(quotes) > 0 {
		var (
			quote     = quotes[len(quotes)-1]
			index     = strings.LastIndexByte(input, quote)
			line, col = getLines(index, input)
		)
		return nil, fmt.Errorf("syntax error: unclosed quote '%c' at line %d, column %d", quote, line, col)
	}

	tokens = appendWord(tokens, current)

	return
}
