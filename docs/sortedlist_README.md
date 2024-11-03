# Sorted List

This package implements generic Sorted List, the data of the list should implement the `IDeepCopy` interface.

## Quick Start
```go
package main

import (
  "github.com/architagr/golang_collections/list"
)


type user struct {
	id   int
	name string
}

func (usr *user) Copy() interface{} {
	cpy := new(user)
	cpy.id = usr.id
	cpy.name = usr.name
	return cpy
}
func (usr *user) Equal(val interface{}) bool {
	data, ok := val.(*user)
	if !ok {
		return false
	}
	return data.id == usr.id && data.name == usr.name
}

func main() {
    list := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
    users := []*user{
		&user{
			id:   1,
			name: "test name 1",
		},
		&user{
			id:   2,
			name: "test name 2",
		},
	}
    index := list.Add(users[0])
    fmt.Println("data was added at index", index)
    data, err := list.Get(0)
    if err != nil{
        panic(fmt.Errorf("Error in getting value at a index: %s", err.Error()))
    }
    err := list.Set(0, &user{
        id: 3,
        name: "test user 3",
    })
	if err != nil {
		 panic(fmt.Errorf("Error in setting value at a index: %s", err.Error()))
	}
    index, err := list.Remove(users[1])
	if err != nil {
		panic(fmt.Errorf("error when removing a node: %s", err.Error()))
	}

    data, err := list.RemoveAtIndex(1)
	if err != nil {
		panic(fmt.Errorf("error when removing a node at a index: %s", err.Error()))
	}

    index := list.Find(&user{
			id:   2,
			name: "test name 2",
		})
    if index == -1 {
		panic(fmt.Errorf("error when finding value: %s", err.Error()))
	}
    response := list.Filter(func(val *user) bool {
		return val.id == 1
	})
}
```
## Functions available

the package exposes below listed functions

### InitSortedList[deepCopy IDeepCopy[T], T any]

created a new sorted list that can have nodes that can hold data of type `IDeepCopy`.
T can be of any data type that implements `IDeepCopy`.
also we have to give a less function that will be used to compare data to maintain sorted order.

### Method in the object of single linked list
#### Add(data *deepCopy) (resultIndex int)

this is a function adds data to the end of the list.

#### Remove(data deepCopy) (removedIndex int, err error)

This function removes data from the list and returns error if the data is not found.

#### RemoveAtIndex(index int) (data *deepCopy, err error)

This function removes data at the index, if the index is not valid then returns error.

#### Count() int

This function return the count of elements in the list.

#### Get(index int) (data deepCopy, err error)

This function gets data at the index, if index is not valid then it returns error.

#### Set(index int, data deepCopy) error

This function updates the data at the index, if index is not valid then returns error

#### Find(data deepCopy) (index int)

This function helps to get the index of first occourance if the data that matches input data and returns index, if index is -1 then the data is not found

#### Filter(f Filterfunc[deepCopy, T]) []deepCopy

This function helps to finding all elements for which the returns true for f, and returns all these elements.

#### RemoveAll(f Filterfunc[deepCopy, T]) []deepCopy

This function helps to remove all elements for which the returns true for f, and also returns all removed elements.