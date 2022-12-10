package queue

import "fmt"

type node[T any] struct {
	data T
	next *node[T]
}

func newNode[T any](value T) *node[T] {
	newNode := new(node[T])
	newNode.data = value
	return newNode
}

type queue[T any] struct {
	count      int
	head, tail *node[T]
}

// IQueue is an interface to be implemented by any queue struct
// helps to know all functions available
type IQueue[T any] interface {
	// Enqueue is used to add new data at tail of the queue
	Enqueue(value T)
	// Dequeue returns error if queue is empty, else returns data at the head of queue
	// and removes this data from head of the queue, practically removing the first node that was added in the queue
	Dequeue() (value T, err error)
	// Peep returns error if queue is empty, else returns data at the head of queue
	// practically returning the first node that was added in the queue
	Peep() (value T, err error)
	// IsEmpty return true if the queue is empty, else returns false
	IsEmpty() bool
	// Count returns the number od nodes present in the queue at given point
	Count() int
}

// NewQueue function helps to initilize a queue of type IQueue
// T is the type of data that you want to have in each node
func NewQueue[T any]() IQueue[T] {
	return new(queue[T])
}
func (queue *queue[T]) Enqueue(value T) {
	node := newNode(value)
	if queue.IsEmpty() {
		queue.head = node
	} else {
		queue.tail.next = node
	}
	queue.tail = node
	queue.count++
}

func (queue *queue[T]) Dequeue() (value T, err error) {
	if queue.IsEmpty() {
		err = fmt.Errorf("queue is empty")
		return
	}
	node := queue.head
	queue.head = queue.head.next
	queue.count--
	if queue.IsEmpty() {
		queue.head = nil
		queue.tail = nil
	}
	value = node.data
	return
}
func (queue *queue[T]) Peep() (value T, err error) {
	if queue.IsEmpty() {
		err = fmt.Errorf("queue is empty")
		return
	}
	value = queue.head.data
	return
}

func (queue *queue[T]) IsEmpty() bool {
	return queue.count == 0
}
func (queue *queue[T]) Count() int {
	return queue.count
}
