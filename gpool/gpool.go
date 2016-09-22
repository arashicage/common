package gpool

import (
	"sync"
)

// Pool struct
type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

// New a pool
func New(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add ...
func (p *Pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

// Done ...
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

// Wait ...
func (p *Pool) Wait() {
	p.wg.Wait()
}
