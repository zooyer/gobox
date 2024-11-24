package gobox

import (
	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/box/bool"
	"github.com/zooyer/gobox/box/cat"
	"github.com/zooyer/gobox/box/echo"
)

var main = map[string]box.Main{
	"cat":   cat.Cat,
	"echo":  echo.Echo,
	"true":  bool.True,
	"false": bool.False,
}

func Get(name string) box.Main {
	return main[name]
}
