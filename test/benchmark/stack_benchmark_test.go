package benchmark

import (
	"testing"

	"github.com/architagr/golang_collections/stack"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
)

func Benchmark_llStack_golang_collections(b *testing.B) {
	newStack := stack.NewStack[int]()
	for i := 0; i < b.N; i++ {
		newStack.Push(i)
	}
}

func Benchmark_llStack_emirpasicGods(b *testing.B) {
	newStack := lls.New()
	for i := 0; i < b.N; i++ {
		newStack.Push(i)
	}
}
