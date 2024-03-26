package main

import (
	"fmt"
	"sync"
)

// 定义一个通用的Pool接口
type Pool interface {
	Get() interface{}
	Put(x interface{})
}

// 使用sync.Pool实现Pool接口
type SyncPool[T any] struct {
	pool *sync.Pool
}

// NewSyncPool 创建一个新的SyncPool
func NewSyncPool[T any]() *SyncPool[T] {
	return &SyncPool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(T) // 这里返回一个新的T类型实例
			},
		},
	}
}

// Get 从池中获取一个对象
func (p *SyncPool[T]) Get() *T {
	v := p.pool.Get().(*T)
	return v
}

// Put 将一个对象放回池中
func (p *SyncPool[T]) Put(x *T) {
	// 这里可以根据需要重置或清理对象x
	p.pool.Put(x)
}

func main() {
	// 创建一个字符串类型的对象池
	stringPool := NewSyncPool[string]()

	// 从池中获取对象
	str := stringPool.Get()
	fmt.Println(str)

	// str = string("Hello, Pool!")

	// 使用完毕后将对象放回池中
	stringPool.Put(str)

	// 再次从池中获取对象，可能会得到之前放回的对象
	str2 := stringPool.Get()
	fmt.Println(str2) // 输出可能是"Hello, Pool!"，也可能是新创建的对象

	// 创建一个整数类型的对象池
	intPool := NewSyncPool[int]()

	// 使用整数类型的对象池...
	num := intPool.Get()
	fmt.Println(num)

	// num = 42

	intPool.Put(num)
}
