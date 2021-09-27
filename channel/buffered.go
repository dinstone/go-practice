package channel

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Pattern() {
	var wg sync.WaitGroup
	wg.Add(4)

	tasks := make(chan string, 10)
	for i := 1; i <= 4; i++ {
		go work(&wg, tasks, i)
	}

	for j := 0; j < 10; j++ {
		tasks <- fmt.Sprintf("Task-%d", j)
	}

	// notify worker close
	close(tasks)

	wg.Wait()
}

func work(wg *sync.WaitGroup, tasks chan string, w int) {
	defer wg.Done()

	for {
		t, ok := <-tasks
		if !ok {
			fmt.Printf("worker %d shutting down\n", w)
			return
		}

		// do work
		fmt.Printf("worker %d started %s\n", w, t)
		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		fmt.Printf("worker %d completed %s\n", w, t)
	}
}
