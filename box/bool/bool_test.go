package bool

import "testing"

func TestNew(t *testing.T) {
	for i := -10; i < 10; i++ {
		if New(i).Main(nil) != i {
			t.Fatal(i)
		}
	}
}
