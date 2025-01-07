package parallel

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestLong(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := longRunningOperation(ctx); err != nil {
		log.Printf("操作被取消: %v", err)
	}
}

func TestSelect(t *testing.T) {
	selectOperation()
}
