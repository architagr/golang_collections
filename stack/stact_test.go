package stack

import (
	"reflect"
	"testing"
)

func TestNewNodeCreatedOfIntType(t *testing.T) {
	node := newNode(10)
	if node == nil {
		t.Errorf("No node was created")
	} else if reflect.ValueOf(node.data).Kind() != reflect.Int {
		t.Errorf("new node was not created with int data type, it was created with %T type", node.data)
	}
}

func TestNewNodeCreatedOfStringType(t *testing.T) {
	node := newNode("test data")
	if node == nil {
		t.Errorf("No node was created")
	} else if reflect.ValueOf(node.data).Kind() != reflect.String {
		t.Errorf("new node was not created with string data type, it was created with %T type", node.data)
	}
}

func TestIntStackToGetCreated(t *testing.T) {
	newStack := NewStack[int]()
	if newStack == nil {
		t.Errorf("No stack was created")
	}
	if newStack.Count() != 0 {
		t.Errorf("wrong initilization of stack with count not equal to zero")
	}
}

func TestStackFunctions(t *testing.T) {
	newStack := NewStack[int]()
	_, err := newStack.Pop()
	if err == nil {
		t.Errorf("should have thrown error as stack does not gave any data")
	}

	newStack.Push(10)
	if newStack.IsEmpty() || newStack.Count() != 1 {
		t.Errorf("after pushing only 1 element to stack we expect the count to be 1")
	}
	newStack.Push(20)
	if newStack.IsEmpty() || newStack.Count() != 2 {
		t.Errorf("after pushing only 1 element to stack we expect the count to be 2")
	}
	data, err := newStack.Pop()
	if err != nil {
		t.Errorf("should be able to get the last value from the stack but received error %s", err.Error())
	}
	if data != 20 {
		t.Errorf("expected 20 to be poped, but got %d", data)
	}

	data = newStack.Top()
	if data != 10 {
		t.Errorf("expected 10 to be on top, but got %d", data)
	}
	if newStack.Count() != 1 {
		t.Errorf("stack should only have 1 value")
	}
	newStack.Pop()
	if !newStack.IsEmpty() {
		t.Errorf("Stack should be empty")
	}
}
