package cmd

import (
	"maps"

	"github.com/zooyer/gobox/box/bool"
	"github.com/zooyer/gobox/box/cat"
	"github.com/zooyer/gobox/box/echo"
	"github.com/zooyer/gobox/box/pwd"
	"github.com/zooyer/gobox/types"
)

var cmd = map[string]types.NewFunc{
	"cat":   cat.New,
	"echo":  echo.New,
	"true":  bool.True,
	"false": bool.False,
	"pwd":   pwd.New,
}

func New(name string) types.NewFunc {
	return cmd[name]
}

func Cmd() map[string]types.NewFunc {
	return maps.Clone(cmd)
}
