package shell

import (
	"errors"
	"flag"
	"fmt"
	"reflect"

	"github.com/zooyer/gobox/types"
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

func writeError(opt types.Option, err error) {
	_, _ = fmt.Fprintln(opt.Stderr, "shell:", err)
}

func deferClose(err *error, close func() error) {
	if close == nil {
		return
	}

	var e = close()
	if err != nil && *err == nil {
		*err = e
	}
}
