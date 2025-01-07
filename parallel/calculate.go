package parallel

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func MasterFunction() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Using %d CPU cores\n", runtime.NumCPU())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	inputChan := make(chan [2]int)
	resultChan := make(chan string)

	wg.Add(1)
	go InputHandle(ctx, inputChan, &wg)

	wg.Add(1)
	go Calculate(ctx, inputChan, resultChan, &wg)

	wg.Add(1)
	go ResultHandle(ctx, resultChan, &wg)

	wg.Wait()
	fmt.Println("All tasks compleated")
}

func ResultHandle(ctx context.Context, resultChan chan string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	count := 0
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				fmt.Println("ResultHandle stopped: result close")
				return
			}
			fmt.Println("Result :", result)
		case <-ctx.Done():
			fmt.Println("ResultHandle stopped: timeout reached")
			return
		}

		count++
		fmt.Printf("ResultHandle loop : %d \n", count)
	}
}

func Calculate(ctx context.Context, inputChan chan [2]int, resultChan chan string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	// defer close(resultChan)

	for {
		select {
		case input, ok := <-inputChan:
			if !ok {
				fmt.Println("Calculate stopped: input close")
				return
			}
			a, b := input[0], input[1]
			add := a + b
			sub := a - b
			result := fmt.Sprintf("Input: %d, %d ; Add: %d, Sub: %d", a, b, add, sub)

			select {
			case resultChan <- result:
			case <-ctx.Done():
				fmt.Println("Calculate stopped: timeout reached")
				return
			}

		case <-ctx.Done():
			fmt.Println("Calculate stopped: timeout reached")
			return
		}
	}

}

func InputHandle(ctx context.Context, inputChan chan [2]int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	defer close(inputChan)

	inputs := [][2]int{
		{4, 2}, {6, 3}, {9, 5}, {18, 7},
	}
	for _, input := range inputs {
		select {
		case inputChan <- input:
			fmt.Printf("Sent input %v\n", input)
		case <-ctx.Done():
			fmt.Println("InputHandle stopped: timeout reached")
			return
		}
	}
	fmt.Println("InputHandle stopped: input finish")
}
