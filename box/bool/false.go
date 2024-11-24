package bool

import (
	"context"

	"github.com/zooyer/gobox/box"
)

func False(ctx context.Context, opt box.Option) (errno int) {
	return 1
}
