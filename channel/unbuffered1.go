package channel

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Court() {
	var wg sync.WaitGroup
	wg.Add(2)

	court := make(chan int)
	go player("townadal", court, &wg)
	go player("dinstone", court, &wg)

	// kick off
	court <- 1

	wg.Wait()
}

func player(pn string, court chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		ball, ok := <-court
		if !ok {
			fmt.Println("player " + pn + " won!")
			return
		}

		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("player %s missed \n", pn)
			close(court)
			return
		}

		fmt.Printf("player %s hit %d\n", pn, ball)
		ball++
		// tran ball to peer
		court <- ball
	}
}
