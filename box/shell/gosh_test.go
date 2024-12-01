package shell

import (
	"context"
	"os"
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestNewGosh(t *testing.T) {
	var option = types.Option{
		Dir:    "",
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	var (
		sh   = NewGosh(option)
		err  error
		code int
	)

	if err = os.Chdir("./testdata"); err != nil {
		t.Fatal(err)
	}

	os.Remove("logs.txt")
	os.Remove("count.txt")
	os.Remove("output.log")

	//sh.Builtin = gobox.Mains()
	sh.Builtin = nil
	//sh.Command = nil

	var commands []Command
	for _, command := range commands {
		if code, err = sh.Exec(&command, option); err != nil {
			t.Error(err)
		}
		if code != 0 {
			t.Fatal(code)
		}
	}
}

func TestGosh(t *testing.T) {
	var option = types.Option{
		Dir:    "",
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	var (
		sh   = NewGosh(option)
		ctx  = context.Background()
		err  error
		code int

		lexer    = NewLexer(option.Stdin)
		parser   = NewParser(lexer.Token())
		errors   = make(chan error, 2)
		commands = parser.Command()
	)

	go func() { errors <- lexer.Run(ctx) }()
	go func() { errors <- parser.Run(ctx) }()

	for command := range commands {
		if code, err = sh.Exec(&command, option); err != nil {
			t.Error(err)
		}
		t.Log("code:", code)
	}

	for err = range errors {
		if err != nil {
			t.Error(err)
		}
	}
}
