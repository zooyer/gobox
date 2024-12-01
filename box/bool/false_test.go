package bool

import (
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestFalse(t *testing.T) {
	if False(types.Option{}).Main(nil) != 1 {
		t.Fatal("false failed")
	}
}
