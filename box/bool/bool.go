package bool

import (
	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/types"
)

type Bool struct {
	box.Process
	errno int
}

func (b *Bool) Main(args []string) (errno int) {
	return b.errno
}

func New(errno int, option types.Option) types.Process {
	return &Bool{
		Process: box.Process{
			Option: option,
		},
		errno: errno,
	}
}
