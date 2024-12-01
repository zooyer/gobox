package box

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

type TestProcess struct {
	Process
}

func (p *TestProcess) Main() {
	p.Start()
	defer p.Stop()

	var ch = make(chan os.Signal, 10)
	defer close(ch)

	p.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range ch {
			fmt.Printf("Got signal: %v\n", sig)
		}
	}()

	for i := 0; i < 10 && p.Run(); i++ {
		fmt.Println("i =", i)
		time.Sleep(1 * time.Second)
	}
}

func TestProcess_Run(t *testing.T) {
	var p TestProcess

	go func() {
		p.Wait()
		fmt.Println("wait done.")
	}()

	go func() {
		time.Sleep(1 * time.Second)
		p.Signal(syscall.SIGINT)
		time.Sleep(2 * time.Second)
		p.Signal(syscall.SIGTERM)
		time.Sleep(2 * time.Second)
		p.Signal(syscall.SIGSTOP)
		time.Sleep(2 * time.Second)
		p.Kill()
		p.Kill()
	}()

	p.Main()
	time.Sleep(1 * time.Second)
}
