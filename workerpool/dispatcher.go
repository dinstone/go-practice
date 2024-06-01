// dispatcher.go

package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ReqHandler is a map of request handlers, keyed by request type.
var ReqHandler = map[int]RequestHandler{
	1: func(data interface{}) error {
		return nil
	},
}

// dispatcher manages a pool of workers and distributes incoming requests among them.
type dispatcher struct {
	inCh        chan Request
	wg          *sync.WaitGroup
	mu          sync.Mutex
	workerCount int
	stopCh      chan struct{} // Channel to signal workers to stop
}

// AddWorker adds a new worker to the pool and increments the worker count.
func (d *dispatcher) AddWorker(w WorkerLauncher) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.workerCount++
	d.wg.Add(1)
	w.LaunchWorker(d.inCh, d.stopCh)
}

// RemoveWorker removes a worker from the pool if the worker count is greater than minWorkers.
func (d *dispatcher) RemoveWorker(minWorkers int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.workerCount > minWorkers {
		d.workerCount--
		d.stopCh <- struct{}{} // Signal a worker to stop
	}
}

// ScaleWorkers dynamically adjusts the number of workers based on the load.
func (d *dispatcher) ScaleWorkers(minWorkers, maxWorkers, loadThreshold int) {
	ticker := time.NewTicker(time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		load := len(d.inCh) // Current load is the number of pending requests in the channel
		if load > loadThreshold && d.workerCount < maxWorkers {
			fmt.Println("Scaling Triggered")
			newWorker := &Worker{
				Wg:         d.wg,
				Id:         d.workerCount,
				ReqHandler: ReqHandler,
			}
			d.AddWorker(newWorker)
		} else if load < int(float32(loadThreshold)*0.75) && d.workerCount > minWorkers {
			fmt.Println("Reducing Triggered")
			d.RemoveWorker(minWorkers)
		}
	}
}

// LaunchWorker launches a worker and increments the worker count.
func (d *dispatcher) LaunchWorker(id int, w WorkerLauncher) {
	w.LaunchWorker(d.inCh, d.stopCh) // Pass stopCh to the worker
	d.mu.Lock()
	d.workerCount++
	d.mu.Unlock()
}

// MakeRequest adds a request to the input channel, or drops it if the channel is full.
func (d *dispatcher) MakeRequest(r Request) {
	select {
	case d.inCh <- r:
	default:
		// Handle the case when the channel is full
		fmt.Println("Request channel is full. Dropping request.")
		// Alternatively, you can log, buffer the request, or take other actions
	}
}

// Stop gracefully stops all workers, waiting for them to finish processing.
func (d *dispatcher) Stop(ctx context.Context) {
	fmt.Println("\nstop called")
	close(d.inCh) // Close the input channel to signal no more requests will be sent
	done := make(chan struct{})

	go func() {
		d.wg.Wait() // Wait for all workers to finish
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All workers stopped gracefully")
	case <-ctx.Done():
		fmt.Println("Timeout reached, forcing shutdown")
		// Forcefully stop all workers if timeout is reached
		for i := 0; i < d.workerCount; i++ {
			d.stopCh <- struct{}{}
		}
	}

	d.wg.Wait()
}

// NewDispatcher creates a new dispatcher with a buffered channel and a wait group.
func NewDispatcher(b int, wg *sync.WaitGroup, maxWorkers int) Dispatcher {
	return &dispatcher{
		inCh:   make(chan Request, b),
		wg:     wg,
		stopCh: make(chan struct{}, maxWorkers), // Buffered channel to prevent blocking on stop
	}
}
