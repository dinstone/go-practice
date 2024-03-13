package pipeline

import (
	"fmt"
	"testing"
	"time"
)

func TestSquaring(t *testing.T) {
	c := gen(2, 3)
	d := squ(c)

	fmt.Println(<-d)
	fmt.Println(<-d)
	fmt.Println(<-d)
}

func TestSquLoop(t *testing.T) {
	for v := range squ(squ(gen(2, 3))) {
		fmt.Println(v)
	}
}

func TestMerge(t *testing.T) {
	c := gen(2, 3)
	d1 := squ(c)
	d2 := squ(c)

	for v := range merge(d1, d2) {
		fmt.Println(v)
	}
}

func TestMerge2(t *testing.T) {

	// Distribute the sq work across two goroutines that both read from in.
	c1 := squ(gen(2, 3))
	c2 := squ(gen(4, 5))

	// Set up a done channel that's shared by the whole pipeline,
	// and close that channel when this pipeline exits, as a signal
	// for all the goroutines we started to exit.
	done := make(chan struct{})
	// defer close(done)

	out := merge2(done, c1, c2)
	// Consume the first value from output.
	fmt.Println(<-out) // 4 or 9
	fmt.Println(<-out)
	fmt.Println(<-out)

	fmt.Println("close done")
	close(done)

}

func AsyncCall(t int) <-chan int {
	c := make(chan int, 1)
	go func() {
		// simulate real task
		time.Sleep(time.Millisecond * time.Duration(t))
		c <- t
	}()
	return c
}

func AsyncCall2(t int) <-chan int {
	c := make(chan int, 1)
	go func() {
		// simulate real task
		time.Sleep(time.Millisecond * time.Duration(t))
		c <- t
	}()
	// gc or some other reason cost some time
	time.Sleep(200 * time.Millisecond)
	return c
}

func TestSelect(t *testing.T) {
	select {
	case resp := <-AsyncCall(50):
		println(resp)
	case resp := <-AsyncCall(200):
		println(resp)
	case resp := <-AsyncCall2(3000):
		println(resp)
	}
}
