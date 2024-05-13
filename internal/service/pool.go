package service

import (
	"sync"
)

type CommandPool struct {
	ch chan func()
}

func NewPool(maxCnt int) *CommandPool {
	pool := &CommandPool{
		ch: make(chan func(), maxCnt),
	}

	go pool.activate()
	return pool
}

func (p *CommandPool) activate() {
	var wg sync.WaitGroup

	for cmd := range p.ch {
		wg.Add(1)
		go func(cmd func()) {
			defer wg.Done()
			cmd()
		}(cmd)
	}

	wg.Wait()
}

func (p *CommandPool) Go(fn func()) {
	p.ch <- fn
}
