package infchan

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	c := NewChan[int](1)

	for i := 0; i < 100; i++ {
		c.Put(i)
	}
	c.Close()
	fmt.Println(c.Len())
	for {

		data, ok := <-c.Get()
		fmt.Println(data, ok, c.Len())
		if !ok {
			break
		}
	}
}
func TestName(t *testing.T) {
	c := NewChan[int](10)
	go func() {
		for i := 0; i < 100; i++ {
			go func() {
				for k := 0; k < 10; k++ {
					c.Put(1)
				}
			}()
		}
	}()

	for i := 0; i < 1000; i++ {
		<-c.Get()
	}

}

func TestCC(t *testing.T) {
	c := make(chan int, 5)
	c <- 1
	c <- 1
	fmt.Println(len(c))
	fmt.Println(cap(c))
}

func BenchmarkChan(b *testing.B) {
	// TODO: Initialize
	c := NewChan[int](1)
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		c.Put(1)
		<-c.Get()
	}
}
