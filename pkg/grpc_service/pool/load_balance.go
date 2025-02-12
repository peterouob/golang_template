package grpcpool

import (
	"sync"
)

type LoadBalance interface {
	Select(int) int
}

type roundRobin struct {
	currIndex int
	mu        sync.Mutex
}

func (rr *roundRobin) Select(n int) int {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	i := rr.currIndex
	rr.currIndex = (i + 1) % n
	return i
}
