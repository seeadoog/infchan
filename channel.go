package infchan

import (
	"container/list"
	"sync"
)

type InfChan[T any] struct {
	data    chan T
	backlog *list.List
	lock    sync.Mutex
	closed  int
}

func NewInfChan[T any](bufferCap int) *InfChan[T] {
	if bufferCap < 1 {
		panic("bufferCap should be > 0")
	}
	return &InfChan[T]{
		data:    make(chan T, bufferCap),
		backlog: list.New(),
	}
}

func (c *InfChan[T]) load() {
	c.lock.Lock()
	if c.backlog.Len() > 0 {
		select {
		case c.data <- c.backlog.Front().Value.(T):
			c.backlog.Remove(c.backlog.Front())
		default:

		}
	} else {
		if c.closed == 1 {
			close(c.data)
			c.closed = 2
		}
	}
	c.lock.Unlock()
}

func (c *InfChan[T]) Get() <-chan T {
	c.load()
	return c.data
}

func (c *InfChan[T]) Put(v T) {
	c.lock.Lock()
	if c.closed > 0 {
		c.lock.Unlock()
		panic("put data to closed channel")
	}
	if c.backlog.Len() == 0 {
		select {
		case c.data <- v:
			c.lock.Unlock()
			return
		default:
		}
	}
	c.backlog.PushBack(v)
	c.lock.Unlock()
}

func (c *InfChan[T]) Len() (length int) {
	c.lock.Lock()
	length = c.backlog.Len() + len(c.data)
	c.lock.Unlock()
	return
}

func (c *InfChan[T]) Close() {
	c.lock.Lock()
	if c.closed > 0 {
		c.lock.Unlock()
		panic("close at closed channel")
	}
	c.closed = 1
	if c.backlog.Len() > 0 {
		c.lock.Unlock()
		return
	}
	c.closed = 2
	c.lock.Unlock()
	close(c.data)

}
