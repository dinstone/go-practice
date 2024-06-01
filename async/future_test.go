package async

import (
	"fmt"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	future := Execute(func() string {
		time.Sleep(time.Second)
		return "hello"
	})

	result, err := future.Get(500)
	if err != nil {
		future.Cancel()
	} else {
		fmt.Println(result)
	}
	//睡眠是为了让任务成功取消，而不是因为程序退出而取消
	time.Sleep(10 * time.Second)
}
