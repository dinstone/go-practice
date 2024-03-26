package pool

import (
	"fmt"
	"sync"
	"time"
)

// Conn 定义了一个通用连接接口
type Conn interface {
	CloseMe()
}

// Factory 定义了一个用于创建新连接的接口
type Factory interface {
	New(*Pool) (Conn, error)
	Close() error
}

// Pool 定义了连接池结构体
type Pool struct {
	lock      sync.Mutex
	conns     map[*Conn]time.Time // 存储当前连接及其最后使用时间
	available chan *Conn          // 空闲可借的队列
	maxActive int                 // 最大活跃连接数
	maxIdle   int                 // 最大空闲连接数
	factory   Factory             // 连接工厂
	minIdle   int                 // 最小空闲连接数
	idleTime  time.Duration       // 空闲连接超时时间
	waitTime  time.Duration       // 等待连接最大时间
	waiting   int                 // 用于监控当前等待连接的请求数
}

// Option 定义了连接池的配置选项
type Option func(*Pool)

// WithMaxActive 设置最大活跃连接数
func WithMaxActive(maxActive int) Option {
	return func(p *Pool) {
		p.maxActive = maxActive
	}
}

// WithMaxIdle 设置最大空闲连接数
func WithMaxIdle(maxIdle int) Option {
	return func(p *Pool) {
		p.maxIdle = maxIdle
	}
}

// WithMinIdle 设置最小空闲连接数
func WithMinIdle(minIdle int) Option {
	return func(p *Pool) {
		p.minIdle = minIdle
	}
}

// WithIdleTimeout 设置空闲连接超时时间
func WithIdleTimeout(idleTime time.Duration) Option {
	return func(p *Pool) {
		p.idleTime = idleTime
	}
}

// WithWaitTime 设置连接最大存活时间
func WithWaitTime(maxLife time.Duration) Option {
	return func(p *Pool) {
		p.waitTime = maxLife
	}
}

// NewPool 创建一个新的连接池
func NewPool(factory Factory, options ...Option) *Pool {
	p := &Pool{
		conns:     make(map[*Conn]time.Time),
		factory:   factory,
		available: make(chan *Conn, 4),
		maxActive: 4,     // 最大活跃连接数
		maxIdle:   2,     // 最大空闲连接数
		minIdle:   1,     // 最小空闲连接数
		idleTime:  10000, // 空闲连接超时时间
		waitTime:  5000,  // 等待连接最大时间
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (cp *Pool) AcquireConn() (*Conn, error) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	// 如果没有空闲连接，检查是否已达到最大活跃连接数
	if cp.maxActive > 0 && len(cp.conns) < cp.maxActive {
		// 创建新连接
		conn, err := cp.factory.New(cp)
		if err != nil {
			return nil, err
		}
		cp.conns[&conn] = time.Now()
		return &conn, nil
	}

	select {
	case conn := <-cp.available:
		cp.waiting-- // 如果有等待的请求，减少等待计数
		return conn, nil
	default:
		cp.waiting++ // 增加等待计数
		// 阻塞，直到有连接可用
		conn, err := cp.waitForConn()
		cp.waiting-- // 减少等待计数
		return conn, err
	}
}

// waitForConn 等待连接变得可用
func (cp *Pool) waitForConn() (*Conn, error) {
	select {
	case conn := <-cp.available:
		return conn, nil
	case <-time.After(cp.waitTime * time.Millisecond): // 设置一个超时，防止永远等待
		return nil, fmt.Errorf("waiting for connection from pool timed out than %dms", cp.waitTime)
	}
}

// ReleaseConn 将连接释放回连接池
func (cp *Pool) ReleaseConn(conn *Conn) error {
	// cp.lock.Lock()
	// defer cp.lock.Unlock()

	if _, ok := cp.conns[conn]; !ok {
		return fmt.Errorf("connection not found in pool")
	}

	cp.available <- conn // 将连接放回可用通道
	return nil
}

// LendingCount 返回当前借出的连接数
func (cp *Pool) LendingCount() int {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	return len(cp.conns) - len(cp.available)
}

// WaitingCount 返回当前等待连接的请求数
func (cp *Pool) WaitingCount() int {
	return cp.waiting
}
