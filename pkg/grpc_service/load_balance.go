package grpc_service

import "sync"

type LoadBalance interface {
	Select(int) int
}

type RoundRobin struct {
	currIndex int
	mu        sync.Mutex
}

var _ LoadBalance = (*RoundRobin)(nil)

func (rr *RoundRobin) Select(n int) int {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	i := rr.currIndex
	rr.currIndex = (i + 1) % n
	return i
}
