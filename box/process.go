package box

import (
	"os"
	"sync"
	"syscall"

	"github.com/zooyer/gobox/types"
)

type Process struct {
	done   bool
	once   sync.Once
	wait   sync.WaitGroup
	mutex  sync.Mutex
	signal map[os.Signal][]chan<- os.Signal
	Option types.Option
}

func (p *Process) Kill() {
	p.Signal(syscall.SIGKILL)
}

func (p *Process) Wait() {
	p.wait.Wait()
}

func (p *Process) Signal(signal os.Signal) {
	switch signal {
	case syscall.SIGSTOP:
		fallthrough
	case syscall.SIGKILL:
		p.Stop()
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, c := range p.signal[signal] {
		c <- signal
	}
}

func (p *Process) Run() bool {
	return !p.done
}

func (p *Process) Stop() {
	p.once.Do(func() {
		p.wait.Done()
		p.done = true
	})
}

func (p *Process) Start() {
	p.wait.Add(1)
}

func (p *Process) Notify(c chan<- os.Signal, sig ...os.Signal) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.signal == nil {
		p.signal = make(map[os.Signal][]chan<- os.Signal)
	}

	for _, s := range sig {
		var has bool
		for _, channel := range p.signal[s] {
			if channel == c {
				has = true
				break
			}
		}

		if !has {
			p.signal[s] = append(p.signal[s], c)
		}
	}
}
