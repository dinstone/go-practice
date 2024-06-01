package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Executor interface {
	Execute(task func(wid int)) error
	Shutdown()
}

type poolExecutor struct {
	coreSize int
	maxSize  int
	queue    chan task
	wg       *sync.WaitGroup
	mu       sync.Mutex
	count    int
	quit     chan struct{} // Channel to signal workers to stop
}

type worker struct {
	id    int
	wg    *sync.WaitGroup
	queue chan task
	quit  chan struct{}
}

func (w *worker) launch() {
	go func() {
		defer w.wg.Done()

		for {
			select {
			case task, open := <-w.queue:
				if !open {
					// If the channel is closed, stop processing and return
					// if we skip close channel check then after closing channel,
					// worker keep reading empty values from closed channel.
					fmt.Println("Stopping worker (queue):", w.id)
					return
				}

				task.run(w.id)
				//time.Sleep(1 * time.Microsecond) // Small delay to prevent tight loop
			case <-w.quit:
				fmt.Println("Stopping worker (quit):", w.id)
				return
			}
		}
	}()
}

type task struct {
	run func(wid int)
}

func (p *poolExecutor) scaleWorker() {
	ticker := time.NewTicker(10 * time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		pendSize := len(p.queue) // Current load is the number of pending task in the channel
		threshold := 3 * p.coreSize
		if pendSize > threshold && p.count < p.maxSize {
			fmt.Println("Scaling increment triggered")
			p.incrementWorker()
		} else if pendSize < int(float32(threshold)*0.75) && p.count > p.coreSize {
			fmt.Println("Scaling decrement triggered")
			p.decrementWorker()
		}

		_, ok := <-p.quit
		if !ok {
			break
		}
	}
	fmt.Println("Scaling shutdown triggered")
}

func (p *poolExecutor) incrementWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.count < p.maxSize {
		p.count++
		p.wg.Add(1)
		worker := worker{id: p.count, wg: p.wg, quit: p.quit, queue: p.queue}
		worker.launch()
	}
}

func (p *poolExecutor) decrementWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.count > p.coreSize {
		p.count--
		p.quit <- struct{}{} // Signal a worker to stop
	}
}

// Execute implements Executor.
func (p *poolExecutor) Execute(fn func(wid int)) error {
	task := task{run: fn}
	select {
	case p.queue <- task:
	default:
		// Handle the case when the channel is full
		fmt.Println("Task channel is full. Dropping task.")
		// Alternatively, you can log, buffer the request, or take other actions
		return errors.New("queue is full and dropping task")
	}

	if p.count < p.coreSize {
		fmt.Println("Execute increment triggered")
		p.incrementWorker()
	}
	return nil
}

// Shutdown implements Executor.
func (p *poolExecutor) Shutdown() {
	fmt.Println("Shutdown called")
	close(p.queue) // Close the input channel to signal no more requests will be sent

	done := make(chan struct{})
	go func() {
		p.wg.Wait() // Wait for all workers to finish
		close(done)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-done:
		fmt.Println("All workers stopped gracefully")
	case <-ctx.Done():
		fmt.Println("Timeout reached, forcing shutdown")
		// Forcefully stop all workers if timeout is reached
		for i := 0; i < p.count; i++ {
			p.quit <- struct{}{}
		}
	}
	close(p.quit)
	p.wg.Wait()
}

func NewExecutor(coreSize, maxSize, queueSize int) Executor {
	var queue chan task
	if queueSize > 0 {
		queue = make(chan task, queueSize)
	} else {
		queue = make(chan task)
	}
	var wg sync.WaitGroup
	p := &poolExecutor{
		coreSize: coreSize, maxSize: maxSize, queue: queue,
		wg: &wg, quit: make(chan struct{}, maxSize),
	}

	go p.scaleWorker()
	return p
}
