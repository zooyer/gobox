package bool

import (
	"testing"

	"github.com/zooyer/gobox/types"
)

func TestNew(t *testing.T) {
	for i := -10; i < 10; i++ {
		if New(i, types.Option{}).Main(nil) != i {
			t.Fatal(i)
		}
	}
}
