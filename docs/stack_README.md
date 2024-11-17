# Stack (using Linked list)

This package implements generic stack, we can have data in node of any type we want. 
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

This benchmarking was done against [emirpasic's GODS package (v1.18.1)](https://github.com/emirpasic/gods)

**System configuration used while doing these benchmark**<br />
**goos:** darwin<br />
**goarch:** arm64<br />
**pkg:** github.com/architagr/golang_collections/test/benchmark<br />
**cpu:** Apple M1 Pro<br />


<table>
<thead>
<tr>
<td>Function name</td>
<td>Package name</td>
<td># operation</td>
<td>Time taken per operation</td>
<td></td>
<td></td>
</tr>
</thead>
<tbody>
<tr >

</tr>
<tr>
<td rowspan=2>Push</td>
<td>architagr/golang_collections</td>
<td>23,863,994</td>
<td>46.03 ns/op</td>
<td>16 B/op</td>
<td>1 allocs/op</td>
</tr>
<tr>
<td>emirpasic/gods</td>
<td>20,692,017</td>
<td>61.07 ns/op</td>
<td>31 B/op</td>
<td>1 allocs/op</td>
</tr>

</tbody>
</table>