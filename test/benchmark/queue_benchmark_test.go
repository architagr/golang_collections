package benchmark

import (
	"testing"

	"github.com/architagr/golang_collections/queue"
	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
)

func Benchmark_llQueueEnqueue_golang_collections(b *testing.B) {
	newQueue := queue.NewQueue[int]()
	for i := 0; i < b.N; i++ {
		newQueue.Enqueue(i)
	}
}

func Benchmark_llQueueEnqueue_emirpasicGods(b *testing.B) {
	newQueue := llq.New()
	for i := 0; i < b.N; i++ {
		newQueue.Enqueue(i)
	}
}
