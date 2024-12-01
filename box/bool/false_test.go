package bool

import "testing"

func TestFalse(t *testing.T) {
	if False().Main(nil) != 1 {
		t.Fatal("false failed")
	}
}
