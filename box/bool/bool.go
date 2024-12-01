package bool

import (
	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/types"
)

type Bool struct {
	box.Process
	code int
}

func (b *Bool) Main(args []string) (code int) {
	return b.code
}

func New(code int, option types.Option) types.Process {
	return &Bool{
		Process: box.Process{
			Option: option,
		},
		code: code,
	}
}
