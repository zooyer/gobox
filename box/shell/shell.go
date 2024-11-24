package shell

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/zooyer/regis/agent/cmd/command"
)

type GNUOption struct {
	Debug       bool   `json:"debug"`
	Debugger    bool   `json:"debugger"`
	DumpPo      bool   `json:"dump-po-strings"`
	DumpStrings bool   `json:"dump-strings"`
	Help        bool   `json:"help"`
	InitFile    string `json:"init-file"`
	Login       bool   `json:"login"`
	NoEditing   bool   `json:"noediting"`
	NoProfile   bool   `json:"noprofile"`
	NoRC        bool   `json:"norc"`
	Posix       bool   `json:"posix"`
	Protected   bool   `json:"protected"`
	RCFile      string `json:"rcfile"`
	Restricted  bool   `json:"restricted"`
	Verbose     bool   `json:"verbose"`
	Version     bool   `json:"version"`
	WordExp     bool   `json:"wordexp"`
}

type RunOption struct {
	Interactive bool `json:"i"`
	Login       bool `json:"l"`
	Restricted  bool `json:"r"`
	Verbose     bool `json:"v"`
	NoClobber   bool `json:"C"`
	Debug       bool `json:"D"`
	NoEditing   bool `json:"n"`
}

type ConfigOption struct {
	AllExport   bool   `json:"a"`
	BraceExpand bool   `json:"B"`
	EmacsEdit   bool   `json:"e"`
	NoBuiltin   bool   `json:"b"`
	ShellFile   string `json:"c"`
	NoProfile   bool   `json:"P"`
}

type Option struct {
	GNUOption
	RunOption
	ConfigOption
}

func bindOption(set *flag.FlagSet, v any) (err error) {
	var pointer = reflect.ValueOf(v)
	if pointer.Kind() != reflect.Ptr || pointer.Elem().Kind() != reflect.Struct {
		return errors.New("must be a pointer to a struct")
	}

	var (
		val = pointer.Elem()
		typ = val.Type()
	)

	for i := 0; i < val.NumField(); i++ {
		var (
			fieldType   = typ.Field(i)
			fieldValue  = val.Field(i)
			jsonTagName = fieldType.Tag.Get("json")
		)

		if fieldValue.Kind() != reflect.Struct && jsonTagName == "" {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Struct:
			if err = bindOption(set, fieldValue.Addr().Interface()); err != nil {
				return
			}
		case reflect.Bool:
			set.BoolVar(fieldValue.Addr().Interface().(*bool), jsonTagName, false, fmt.Sprintf("Set %s flag", jsonTagName))
		case reflect.String:
			set.StringVar(fieldValue.Addr().Interface().(*string), jsonTagName, "", fmt.Sprintf("Set %s option", jsonTagName))
		default:
			return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}
	}

	return
}

func writeError(attr command.Attr, err error) {
	_, _ = fmt.Fprintln(attr.Stderr, "shell:", err)
}

// 判断是否需要继续读取
func needsContinuation(end, input string) (string, bool) {
	input = strings.TrimSpace(input)

	var index int
	if end != "" {
		if index = strings.Index(input, end); index == -1 {
			return end, true
		}

		input = input[index:]
	}

	if len(input) == 0 {
		return "", false
	}

	for len(input) > 0 {
		if end != "" {
			if index = strings.Index(input, end); index == -1 {
				return end, true
			}

			end = ""
			input = input[index:]
		}

		for begin, unit := range beginUints {
			if index = strings.Index(input, begin); index < 0 {
				continue
			}

			if index += len(begin); index < len(input) {
				input = input[index:]
			} else {
				input = input[:]
			}

			_, end = unit.Token()

			break
		}
	}

	return end, true
}

func parseCommand() {

}
