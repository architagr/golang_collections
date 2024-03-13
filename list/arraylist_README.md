# Array List

This package implements generic Array List, the data of the list should implement the `IDeepCopy` interface.

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
    list := InitArrayList[user, *user]()
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
    err := obj.AddAtIndex(0, users[0])
    if err != nil{
        panic(fmt.Errorf("Error in adding at 0th index: %s", err.Error()))
    }
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

    index := list.Find(func(val *user) bool {
		return val.id == 1
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

### InitArrayList[T any, deepCopy IDeepCopy[T]]

created a new array list that can have nodes that can hold data of type `IDeepCopy`.
T can be of any data type that implements `IDeepCopy`.

### Method in the object of single linked list
#### Add(data *deepCopy) (resultIndex int)

this is a function adds data to the end of the list.

#### AddAtIndex(index int, data deepCopy) (err error)

this function adds data to the index in the list and shifts all data to right if the index is out of bound then return error.

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

#### Find(f filterfunc[T, deepCopy]) (index int)

This function helps to get the first occourance if the data that matches according to the filter func and also returns index, if index is -1 then the data is not found

