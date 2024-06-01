package pool

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestExecutor(t *testing.T) {
	// Set GOMAXPROCS to the number of available CPUs
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("Running with %d CPUs\n", numCPU)

	executor := NewExecutor(3, 20, 40)
	for i := 0; i < 50; i++ {
		index := i
		executor.Execute(func(wid int) {
			//time.Sleep(10 * time.Millisecond)
			fmt.Printf("(worker: %d) => (task: %d) -> ok\n", wid, index)
		})
	}

	fmt.Println("Exiting test")

	time.Sleep(15 * time.Second)
	executor.Shutdown()
}
