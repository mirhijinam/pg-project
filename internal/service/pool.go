package service

import (
	"sync"
)

type Pool struct {
	ch chan func()
}

func NewPool(maxCnt int) *Pool {
	pool := &Pool{
		ch: make(chan func(), maxCnt),
	}

	go pool.activate()
	return pool
}

func (p *Pool) activate() {
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

func (p *Pool) Go(fn func()) {
	p.ch <- fn
}
