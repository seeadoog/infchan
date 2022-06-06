package infchan

import (
	"container/list"
	"sync"
)

type InfChan struct {
	data    chan interface{}
	backlog *list.List
	lock    sync.Mutex
	closed  int
}

func NewInfChan(bufferCap int) *InfChan {
	if bufferCap < 1 {
		panic("bufferCap should be > 0")
	}
	return &InfChan{
		data:    make(chan interface{}, bufferCap),
		backlog: list.New(),
	}
}

func (c *InfChan) load() {
	c.lock.Lock()
	if c.backlog.Len() > 0 {
		select {
		case c.data <- c.backlog.Front().Value:
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

func (c *InfChan) Get() <-chan interface{} {
	c.load()
	return c.data
}

func (c *InfChan) Put(v interface{}) {
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

func (c *InfChan) Len() (length int) {
	c.lock.Lock()
	length = c.backlog.Len() + len(c.data)
	c.lock.Unlock()
	return
}

func (c *InfChan) Close() {
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
