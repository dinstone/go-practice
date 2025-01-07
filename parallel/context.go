package parallel

import (
	"context"
	"fmt"
	"log"
	"time"
)

func longRunningOperation(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 执行一些工作
			log.Printf("default running")
			time.Sleep(7 * time.Second)
		}
	}
}

func selectOperation() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(20 * time.Second)
		ch1 <- 1
		fmt.Println("ch1 ending")
	}()

	go func() {
		time.Sleep(25 * time.Second)
		ch2 <- 2
		fmt.Println("ch2 ending")
	}()

	select {
	case data := <-ch1:
		fmt.Println("收到ch1的数据:", data)
	case data := <-ch2:
		fmt.Println("收到ch2的数据:", data)
	case <-time.After(60 * time.Second):
		fmt.Println("超时：没有接收到任何数据")
	}
}
