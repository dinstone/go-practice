package channel

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Baton() {
	var wg sync.WaitGroup
	wg.Add(1)

	baton := make(chan int)

	go runner(&wg, baton)

	// kick off
	baton <- 1

	wg.Wait()
}

func runner(s *sync.WaitGroup, baton chan int) {
	run := <-baton
	fmt.Printf("runner %d running with baton\n", run)

	//next run
	next := run + 1
	if run != 4 {
		fmt.Printf("runner %d to be the line\n", next)
		go runner(s, baton)
	}

	// running
	time.Sleep(100 * time.Millisecond)

	if run == 4 {
		fmt.Printf("runner %d finished, Race over\n", run)
		s.Done()
		return
	}

	fmt.Printf("runner %d exchange with runner %d\n", run, next)
	baton <- next
}
