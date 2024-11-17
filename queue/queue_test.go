package queue

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

func TestIntQueueToGetCreated(t *testing.T) {
	newQueue := NewQueue[int]()
	if newQueue == nil {
		t.Errorf("No queue was created")
	}
	if newQueue.Count() != 0 {
		t.Errorf("wrong initilization of queue with count not equal to zero")
	}
}

func TestStackFunctions(t *testing.T) {
	newQueue := NewQueue[int]()
	_, err := newQueue.Peep()
	if err == nil {
		t.Errorf("should have thrown error as queue does not have any data")
	}
	_, err = newQueue.Dequeue()
	if err == nil {
		t.Errorf("should have thrown error as queue does not have any data")
	}

	newQueue.Enqueue(10)
	if newQueue.IsEmpty() || newQueue.Count() != 1 {
		t.Errorf("after pushing only 1 element to queue we expect the count to be 1")
	}
	newQueue.Enqueue(20)
	if newQueue.IsEmpty() || newQueue.Count() != 2 {
		t.Errorf("after pushing only 1 element to queue we expect the count to be 2")
	}
	data, err := newQueue.Dequeue()
	if err != nil {
		t.Errorf("should be able to get the first value from the queue but received error %s", err.Error())
	}
	if data != 10 {
		t.Errorf("expected 10 to be Dequeueed, but got %d", data)
	}

	data, err = newQueue.Peep()
	if err != nil {
		t.Errorf("should be able to get the first value from the queue but received error %s", err.Error())
	}
	if data != 20 {
		t.Errorf("expected 20 to be on top, but got %d", data)
	}
	if newQueue.Count() != 1 {
		t.Errorf("queue should only have 1 value")
	}
	newQueue.Dequeue()
	if !newQueue.IsEmpty() {
		t.Errorf("queue should be empty")
	}
}
