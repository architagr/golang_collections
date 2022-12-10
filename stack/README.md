# Stack (using Linked list)

This package implements genric stack, we can have data in node of any type we want. 
Stack is a data structure that follows LIFO (Last in first out).

Examples of stacks are: 
1. Pile of books
2. Clipboards in the computers (ctrl + c and ctrl + v)


## Quick Start
```go
package main

import (

  "github.com/architagr/golang_collections/stack"
)

func main() {
  newStack := NewStack[int]()
  newStack.Push(10)
  newStack.Push(11)
  data, err := newStack.Pop()
  if err != nil{
    panic(fmt.Errorf("Error in pop: %s", err.Error()))
  }
  fmt.Printf("data poped %d\n", data)

  data, err = newStack.Top()
  if err != nil{
    panic(fmt.Errorf("Error in top: %s", err.Error()))
  }
  fmt.Printf("data at top %d\n", data)
  fmt.Printf("Stack contains #%d nodes\n", newStack.Count())
  fmt.Printf("Is Stack empty? #%t\n", newStack.IsEmpty())
}
```
## Functions available

the package exposes below listed functions

### NewStack[T any]

created a new stack that can have nodes that can hold data of type T.
T can be of any data type from basic data type like int, string, etc. to any user created data type like structs

### Push(value T)

this is a function been implemented by type struct and helps in adding new node with data been passed as argument to this function. the new node is always added on the top (LIFO)

### Pop() (value T, err error)

this function returns error if we do not have any data/node in the stack, i.e. the stack is empty.
if stack is not empty then it removes the last node that was added and returns the data of the node that was removed.

### Top() T

This function returns last value that was added, without removing the value from the stack.

### Count() int

This function returns the number of nodes in the stack

### IsEmpty() bool

This function returns ture if the stack is empty, i.e. count is 0, and return false if we have some data in stack, i.e. count > 0.


## Benchmarks

This benchmarking was done against [emirpasic's GODS package](https://github.com/emirpasic/gods)

**System configuration used while doing these benchmark**<br />
**goos:** linux<br />
**goarch:** amd64<br />
**pkg:** github.com/architagr/golang_collections/stack<br />
**cpu:** Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz<br />

| Benchmark name                       |       (1)     |             (2) |          (3)  |             (4) |
| ------------------------------------ | -------------:| ---------------:| -------------:| ---------------:|
| BenchmarkGolang_collections_llStack  | **100000000** | **109.0 ns/op** |   **16 B/op** | **1 allocs/op** |
| BenchmarkEmirpasicGods_llStack       |     100000000 |     143.0 ns/op |       31 B/op |     1 allocs/op |