package proxy

import (
	"fmt"
	"testing"
	"time"
)

func TestMethod_Invoke(t *testing.T) {
	//这里对活动时长做统计
	people := &People{}   //创建目标类
	h := new(PeopleProxy) //创建接口实现类
	proxy := NewProxy(people, h)
	//调用方法
	ret, err := proxy.InvokeMethod("Work", "打游戏", "学习")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ret)
}

// 目标类
type People struct {
}

func (p *People) Work(content string, next string) string {
	fmt.Println("活动内容是：" + content + "，接下来需要做：" + next)
	return "all right"
}

// 用户需要自己实现的增强内容，需要实现InvocationHandler接口
type PeopleProxy struct {
}

// 在这里做方法增强
func (p *PeopleProxy) Invoke(proxy *Proxy, method *Method, args []interface{}) ([]interface{}, error) {
	start := time.Now()
	defer fmt.Printf("耗时：%v\n", time.Since(start))
	fmt.Println("before method")
	invoke, err := method.Invoke(args...)
	fmt.Println("after method")
	return invoke, err
}
