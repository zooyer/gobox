package bool

import "github.com/zooyer/gobox/types"

func False(opt types.Option) types.Process {
	return New(1, opt)
}
