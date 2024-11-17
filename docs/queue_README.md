# Queue (using Linked list)

This package implements generic queue, we can have data in node of any type we want. 
Queue is a data structure that follows FIFO (First in First out).

Examples of stacks are: 
1. Queue at the ticket counter
2. Patients waiting outside the doctor's clinic


## Quick Start
```go
package main

import (

  "github.com/architagr/golang_collections/queue"
)

func main() {
  newQueue := NewQueue[int]()
  newQueue.Enqueue(10)
  newQueue.Enqueue(11)
  data, err := newQueue.Dequeue()
  if err != nil{
    panic(fmt.Errorf("Error in pop: %s", err.Error()))
  }
  fmt.Printf("data poped %d\n", data)

  data, err = newQueue.Peep()
  if err != nil{
    panic(fmt.Errorf("Error in top: %s", err.Error()))
  }
  fmt.Printf("data at top %d\n", data)
  fmt.Printf("Queue contains #%d nodes\n", newQueue.Count())
  fmt.Printf("Is Queue empty? #%t\n", newQueue.IsEmpty())
}
```
## Functions available

the package exposes below listed functions

### NewQueue[T any]

created a new queue that can have nodes that can hold data of type T.
T can be of any data type from basic data type like int, string, etc. to any user created data type like structs

### Enqueue(value T)

this is a function been implemented by type struct and helps in adding new node with data been passed as argument to this function. The new node is always added on the tail (FIFO)

### Dequeue() (value T, err error)

this function returns error if we do not have any data/node in the queue, i.e. the queue is empty.
if queue is not empty then it removes the first node that was added and returns the data of the node that was removed.

### Top() T

This function returns first value that was added, without removing the value from the queue.

### Count() int

This function returns the number of nodes in the queue

### IsEmpty() bool

This function returns ture if the queue is empty, i.e. count is 0, and return false if we have some data in queue, i.e. count > 0.


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
<td rowspan=2>Enqueue</td>
<td>architagr/golang_collections</td>
<td>30,853,053</td>
<td>43.30 ns/op</td>
<td>16 B/op</td>
<td>1 allocs/op</td>
</tr>
<tr>
<td>emirpasic/gods</td>
<td>19,681,999</td>
<td>64.10 ns/op</td>
<td>31 B/op</td>
<td>1 allocs/op</td>
</tr>

</tbody>
</table>
