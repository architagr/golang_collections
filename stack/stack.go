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
type IStack[T any] interface {
	Push(value T)
	Pop() (value T, err error)
	Top() (value T)
	IsEmpty() bool
	Count() int
}

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
		err = fmt.Errorf("No data to pop")
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
