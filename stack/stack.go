package stack

import "fmt"

type node[T any] struct {
	data T
	prev *node[T]
}

func newNode[T any](value T) *node[T] {
	newNode := new(node[T])
	newNode.data = value
	return newNode
}

type stack[T any] struct {
	count int
	head  *node[T]
}

// IStack is an interface to be implemented by any stack struct
// helps to know all functions available
type IStack[T any] interface {
	// Push is used to push new data in the stack's top
	Push(value T)
	// Pop returns error if stack is empty, or returns the data in the top node
	// and moves top to next node, practically deleting the last node added
	Pop() (value T, err error)
	// Top return data of the last node been addd
	Top() (value T)
	// IsEmpty returns true if there is no data in the stack, else returns false
	IsEmpty() bool
	// Count returns number of nodes in the stack
	Count() int
}

// NewStack function helps to initilize a stack of type IStack
// T is the type of data that you want to have in the node
func NewStack[T any]() IStack[T] {
	return new(stack[T])
}

func (stack *stack[T]) Push(value T) {
	newNode := newNode(value)
	stack.count++
	if stack.head == nil {
		stack.head = newNode
		return
	}
	newNode.prev = stack.head
	stack.head = newNode

}

func (stack *stack[T]) Pop() (value T, err error) {
	if stack.head == nil {
		err = fmt.Errorf("no data to pop")
		return
	}
	node := stack.head
	stack.head = stack.head.prev
	stack.count--
	value = node.data
	return
}

func (stack *stack[T]) Top() (value T) {
	if stack.head == nil {
		return
	}
	value = stack.head.data
	return
}

func (stack *stack[T]) IsEmpty() bool {
	return stack.count == 0
}

func (stack *stack[T]) Count() int {
	return stack.count
}
