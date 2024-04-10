package workpool

import "sync"

func Publish(numbers []int) <-chan int {
	workChan := make(chan int)
	go func() {
		defer close(workChan)

		for _, v := range numbers {
			workChan <- v
		}
	}()
	return workChan
}

func PoolProcess(numbers []int) []*data {
	workChan := Publish(numbers)
	outChan := make(chan *data)

	var wg sync.WaitGroup
	count := 5
	for i := 0; i < count; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for n := range workChan {
				outChan <- &data{
					number: n,
					square: float64(n) * float64(n),
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outChan)
	}()

	outputs := make([]*data, 0)
	for d := range outChan {
		outputs = append(outputs, d)
	}
	return outputs
}
