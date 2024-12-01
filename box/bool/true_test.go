package bool

import (
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestTrue(t *testing.T) {
	if True(types.Option{}).Main(nil) != 0 {
		t.Fatal("true failed")
	}
}
