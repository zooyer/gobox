package test

import (
	"fmt"
	"time"

	"github.com/zooyer/gobox/box"
	"github.com/zooyer/gobox/types"
)

type Test struct {
	box.Process
}

func (t *Test) Main(args []string) (errno int) {
	t.Start()
	defer t.Stop()

	fmt.Println("args:", args)

	for i := 0; i < 10; i++ {
		fmt.Println("i :", i)
		time.Sleep(1 * time.Second)
	}

	return 0
}

func New(option types.Option) types.Process {
	return &Test{
		Process: box.Process{
			Option: option,
		},
	}
}
