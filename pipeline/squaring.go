package pipeline

import (
	"fmt"
	"sync"
)

//
// Go中构建流数据pipeline
// https://mp.weixin.qq.com/s/E0dFzNVISOWyx00Kx0SZAA
//

func gen(nums ...int) <-chan int {
	outChan := make(chan int)
	go func() {
		for _, n := range nums {
			outChan <- n
		}
		close(outChan)
	}()
	return outChan
}

func squ(inChan <-chan int) <-chan int {
	outChan := make(chan int)
	go func() {
		for n := range inChan {
			fmt.Println("n =", n)
			outChan <- n * n
		}
		close(outChan)
	}()
	return outChan
}

func merge(ins ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	// define merge action
	output := func(in <-chan int) {
		for v := range in {
			out <- v
		}
		wg.Done()
	}
	// add wait count
	wg.Add(len(ins))
	// Start an output goroutine for each input
	for _, in := range ins {
		go output(in)
	}
	// Start a goroutine to close out
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func merge2(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed or it receives a value
	// from done, then output calls wg.Done.
	output := func(index int, c <-chan int) {
		defer wg.Done()

		for {
			select {
			case x, ok := <-c:
				if !ok {
					fmt.Println(index, "input event")
					return
				}
				out <- x
			case <-done:
				fmt.Println(index, "done event")
				return
			}
		}

		// for n := range c {
		// 	select {
		// 	case out <- n:
		// 		fmt.Println("out event")
		// 	case <-done:
		// 		fmt.Println("done event")
		// 	}
		// }
		// fmt.Println("waite done")
	}

	// add wait count
	wg.Add(len(cs))
	// Start an output goroutine for each input
	for index, in := range cs {
		go output(index, in)
	}
	// Start a goroutine to close out
	go func() {
		wg.Wait()
		fmt.Println("output closed")
		close(out)
	}()
	return out
}
