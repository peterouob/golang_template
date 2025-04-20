package utils

import (
	"sync"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	const goroutineCount = 100000
	wg := sync.WaitGroup{}
	wg.Add(goroutineCount)

	s := NewIdWorker(100)
	idSet := sync.Map{}

	for range make([]struct{}, goroutineCount) {
		go func() {
			id := s.GenID()

			if _, loaded := idSet.LoadOrStore(id, struct{}{}); loaded {
				t.Errorf("duplicate ID detected: %d", id)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	t.Log("All IDs are unique")
}

func BenchmarkGenID(b *testing.B) {
	s := NewIdWorker(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.GenID()
		}
	})
}
