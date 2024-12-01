package shell

import (
	"context"
	"strings"
)

type Command struct {
	Path       string
	Args       []string
	Pipe       *Command // |
	Or         *Command // ||
	And        *Command // &&
	Input      string   // <
	Output     string   // >
	Append     string   // >>
	Heredoc    string   // <<
	Background bool     // &

	//Next       *Command // ;
	//Child      *Command // $()
}

func (c *Command) CmdArgs() []string {
	return append([]string{c.Path}, c.Args...)
}

type Parser struct {
	tokens   <-chan Token
	commands chan Command
}

func NewParser(tokens <-chan Token) *Parser {
	return &Parser{
		tokens:   tokens,
		commands: make(chan Command, 1024),
	}
}

func (p *Parser) Command() <-chan Command {
	return p.commands
}

func (p *Parser) put(command Command) {
	if command.Path == "" {
		return
	}

	p.commands <- command
}

func (p *Parser) token(ctx context.Context) (token Token, run bool) {
	select {
	case <-ctx.Done():
		return Token{}, false
	case token, run = <-p.tokens:
		return
	}
}

func (p *Parser) Run(ctx context.Context) (err error) {
	defer close(p.commands)

	var (
		run        bool
		token      Token
		command    Command
		current    = &command
		front      = current
		background bool
	)

	for {
		if token, run = p.token(ctx); !run {
			break
		}

		if background {
			if token.Type == TokenHeredoc {
				command.Heredoc = token.Value
			}

			p.put(command)
			command = Command{}
			current = &command

			background = false

			if token.Type == TokenHeredoc {
				continue
			}
		}

		switch token.Type {
		case TokenWord:
			if current.Path == "" {
				current.Path = token.Value
			} else {
				current.Args = append(current.Args, token.Value)
			}
		case TokenPipe:
			current.Pipe = new(Command)
			current = current.Pipe
			front = current
		case TokenOr:
			front.Or = new(Command)
			current = front.Or
			front = current
		case TokenAnd:
			current.And = new(Command)
			current = current.And
		case TokenRedirectIn:
			current.Input = (<-p.tokens).Value
		case TokenRedirectOut:
			current.Output = (<-p.tokens).Value
		case TokenRedirectAppend:
			current.Append = (<-p.tokens).Value
		case TokenHeredoc:
			current.Heredoc = token.Value
		case TokenBackground:
			background = true
			current.Background = true
		case TokenSemicolon:
			p.put(command)
			command = Command{}
			current = &command
		case TokenVar:
		case TokenCmd:
		}
	}

	p.put(command)

	return nil
}

func ParseCommands(input string) (commands []Command, err error) {
	var (
		ctx     = context.Background()
		lexer   = NewLexer(strings.NewReader(input))
		parser  = NewParser(lexer.Token())
		command = parser.Command()
		errors  = make(chan error, 20)
	)

	go func() { errors <- lexer.Run(ctx) }()

	go func() { errors <- parser.Run(ctx) }()

	for c := range command {
		commands = append(commands, c)
	}

	close(errors)

	for err = range errors {
		if err != nil {
			return
		}
	}

	return
}
