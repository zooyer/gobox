package bool

import "testing"

func TestTrue(t *testing.T) {
	if True().Main(nil) != 0 {
		t.Fatal("true failed")
	}
}
