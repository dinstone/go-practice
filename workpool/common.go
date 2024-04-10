package workpool

import "sync"

func Process(numbers []int) []*data {
	outChan := make(chan *data)

	var wg sync.WaitGroup
	for _, v := range numbers {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			outChan <- &data{
				number: n,
				square: float64(n) * float64(n),
			}
		}(v)
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
