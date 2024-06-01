package workerpool

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	// Set GOMAXPROCS to the number of available CPUs
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("Running with %d CPUs\n", numCPU)

	// Configuration
	bufferSize := 50000
	maxWorkers := 20
	minWorkers := 3
	loadThreshold := 40000
	requests := 50000

	var wg sync.WaitGroup
	dispatcher := NewDispatcher(bufferSize, &wg, maxWorkers)

	// Start initial set of workers
	for i := 0; i < minWorkers; i++ {
		fmt.Printf("Starting worker with id %d\n", i)
		w := &Worker{
			Wg:         &wg,
			Id:         i,
			ReqHandler: ReqHandler,
		}
		dispatcher.AddWorker(w)
	}

	// Start the scaling logic in a separate goroutine
	go dispatcher.ScaleWorkers(minWorkers, maxWorkers, loadThreshold)

	// Send requests to the dispatcher
	for i := 0; i < requests; i++ {
		req := Request{
			Data:    fmt.Sprintf("(Msg_id: %d) -> Hello", i),
			Handler: func(result interface{}) error { return nil },
			Type:    1,
			Timeout: 5 * time.Second,
		}
		dispatcher.MakeRequest(req)
	}

	// Gracefully stop the dispatcher
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dispatcher.Stop(ctx)
	fmt.Println("Exiting main!")

	t.Logf("over")
}
