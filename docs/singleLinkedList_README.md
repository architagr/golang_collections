# Single List list

This package implements generic Single linked list, the data of the list should implement the `IDeepCopy` interface.

## Quick Start
```go
package main

import (
  "github.com/architagr/golang_collections/list"
)


type Integer int

func InitInteger(val int) Integer {
	return Integer(val)
}
func (obj Integer) Copy() interface{} {
	return obj
}

func (obj Integer) Equal(val interface{}) bool {
	x := val.(Integer)
	return int(x) == int(obj)
}

func main() {
    linkedList := InitSingleLinkedList[*Integer]()
    intVal := InitInteger(0)
    linkedList.Add(&intVal)
    intVal2 := InitInteger(1)
    err := obj.AddAtIndex(0, &intVal2)
    if err != nil{
        panic(fmt.Errorf("Error in adding at 0th index: %s", err.Error()))
    }
    data, err := linkedList.Get(0)
    if err != nil{
        panic(fmt.Errorf("Error in getting value at a index: %s", err.Error()))
    }
    intVal3 := InitInteger(2)
    err := linkedList.Set(0, &intVal3)
	if err != nil {
		 panic(fmt.Errorf("Error in setting value at a index: %s", err.Error()))
	}
    index, err := linkedList.Remove(&intVal2)
	if err != nil {
		panic(fmt.Errorf("error when removing a node: %s", err.Error()))
	}

    data, err := linkedList.RemoveAtIndex(1)
	if err != nil {
		panic(fmt.Errorf("error when removing a node at a index: %s", err.Error()))
	}
    x := InitInteger(2)
    index := obj.Find(x)
    if index == -1 {
		panic(fmt.Errorf("error when finding value: %s", err.Error()))
	}

}
```
## Functions available

the package exposes below listed functions

### InitSingleLinkedList[deepCopy IDeepCopy[T], T any]

created a new single linked list that can have nodes that can hold data of type `IDeepCopy`.
T can be of any data type that implements `IDeepCopy`.

### Method in the object of single linked list
#### Add(data deepCopy) (resultIndex int)

this is a function adds data to the end of the list.

#### AddAtIndex(index int, data deepCopy) (err error)

this function adds data to the index in the list and shifts all data to right if the index is out of bound then return error.

#### Remove(data deepCopy) (removedIndex int, err error)

This function removes data from the list and returns error if the data is not found.

#### RemoveAtIndex(index int) (data deepCopy, err error)

This function removes data at the index, if the index is not valid then returns error.

#### Count() int

This function return the count of elements in the list.

#### Get(index int) (data deepCopy, err error)

This function gets data at the index, if index is not valid then it returns error.

#### Set(index int, data deepCopy) error

This function updates the data at the index, if index is not valid then returns error

#### Find(data deepCopy) (index int)

This function helps to get the index of first occourance if the data that matches input data and returns index, if index is -1 then the data is not found

#### RemoveAll(f Filterfunc[deepCopy, T]) []deepCopy

This function helps to remove all elements for which the returns true for f, and also returns all removed elements.
