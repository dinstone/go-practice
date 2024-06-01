package async

import (
	"errors"
	"fmt"
	"time"
)

type Future[T any] struct {
	resultChan chan T
	cancelChan chan struct{}
}

func (future *Future[T]) Get(timeout int) (T, error) {
	timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	defer timer.Stop()
	var tmp T
	select {
	case data := <-future.resultChan:
		return data, nil
	case <-timer.C:
		fmt.Println("Get超时")
		return tmp, errors.New("time out")
	}
}

func (future *Future[T]) Cancel() {
	future.cancelChan <- struct{}{}
}

func Execute[T any](task func() T) *Future[T] {
	future := &Future[T]{
		resultChan: make(chan T, 1),
		cancelChan: make(chan struct{}, 1),
	}
	go func() {
		tmp := make(chan T)
		go func() {
			tmp <- task()
		}()
		select {
		case result := <-tmp:
			future.resultChan <- result
			return
		case <-future.cancelChan:
			fmt.Println("任务已取消")
			return
		}
	}()
	return future
}
