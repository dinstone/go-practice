package pattern

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	maxGoroutines = 25
	poolResources = 2
)

type DbConnection struct {
	ID int32
}

func (c *DbConnection) Close() error {
	log.Println("connection closed:", c.ID)
	return nil
}

var idCounter int32

func createConnection() (io.Closer, error) {
	atomic.AddInt32(&idCounter, 1)
	log.Println("create connection:", idCounter)
	return &DbConnection{idCounter}, nil
}

func TestPool(t *testing.T) {
	log.Println("startup ...")
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	p, e := New(createConnection, poolResources)
	if e != nil {
		log.Println(e)
		return
	}

	for i := 0; i < maxGoroutines; i++ {
		go func(q int) {
			doquery(p, q)
			wg.Done()
		}(i)
	}

	wg.Wait()
	p.Close()
	log.Println("shutdown ...")
}

func doquery(p *Pool, q int) {
	c, e := p.Acquire()
	if e != nil {
		log.Println(e)
		return
	}

	// release connection
	defer p.Release(c)

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	// type cast
	log.Printf("QID[%d] CID[%d]\n", q, c.(*DbConnection).ID)
}
