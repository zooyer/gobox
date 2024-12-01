package bool

import "github.com/zooyer/gobox/types"

func True(opt types.Option) types.Process {
	return New(0, opt)
}
