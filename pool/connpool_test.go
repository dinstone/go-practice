package pool

import (
	"fmt"
	"testing"
	"time"
)

type ConnFactory struct {
}

type PooledConn struct {
	pool *Pool
}

func (f *ConnFactory) New(pool *Pool) (Conn, error) {
	return PooledConn{pool: pool}, nil
}

func (f *ConnFactory) Close() error { return nil }

func (pc PooledConn) CloseMe() {
	pc.pool.ReleaseConn(nil)
}

func TestMain(t *testing.T) {
	connFactory := ConnFactory{}
	pool := NewPool(&connFactory, WithMaxActive(4), WithWaitTime(3000))
	for i := 0; i < 10; i++ {
		go func(id int) {
			conn, err := pool.AcquireConn()
			if err != nil {
				fmt.Printf("Goroutine %d: %v\n", id, err)
				return
			}
			fmt.Printf("Goroutine %d: Borrowed connection: %p\n", id, conn)

			// 模拟使用连接进行工作
			time.Sleep(time.Second)

			fmt.Printf("Goroutine %d: awake\n", id)

			// 释放连接回连接池
			err = pool.ReleaseConn(conn)
			if err != nil {
				fmt.Printf("Goroutine %d: %v\n", id, err)
				return
			}
			fmt.Printf("Goroutine %d: Released connection: %p\n", id, conn)
		}(i)
	}

	select {}
}
