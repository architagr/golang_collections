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
    userList := InitArrayList[*user]()
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
    data, err := userList.Get(0)
    if err != nil{
        panic(fmt.Errorf("Error in getting value at a index: %s", err.Error()))
    }
    err := userList.Set(0, &user{
        id: 3,
        name: "test user 3",
    })
	if err != nil {
		 panic(fmt.Errorf("Error in setting value at a index: %s", err.Error()))
	}
    index, err := userList.Remove(users[1])
	if err != nil {
		panic(fmt.Errorf("error when removing a node: %s", err.Error()))
	}

    data, err := userList.RemoveAtIndex(1)
	if err != nil {
		panic(fmt.Errorf("error when removing a node at a index: %s", err.Error()))
	}

    index := userList.Find(&user{
			id:   2,
			name: "test name 2",
		})
    if index == -1 {
		panic(fmt.Errorf("error when finding value: %s", err.Error()))
	}
    response := userList.Filter(func(val *user) bool {
		return val.id == 1
	})
}
```
## Functions available

the package exposes below listed functions

### InitArrayList[deepCopy IDeepCopy[T], T any]

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

#### Find(data deepCopy) (index int)

This function helps to get the index of first occourance if the data that matches input data and returns index, if index is -1 then the data is not found

#### Filter(f Filterfunc[deepCopy, T]) []deepCopy

This function helps to finding all elements for which the returns true for f, and returns all these elements.

#### RemoveAll(f Filterfunc[deepCopy, T]) []deepCopy

This function helps to remove all elements for which the returns true for f, and also returns all removed elements.


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
<td rowspan=2>Add</td>
<td>architagr/golang_collections</td>
<td>74,762,736</td>
<td>21.34 ns/op</td>
<td>22 B/op</td>
<td>0 allocs/op</td>
</tr>
<tr>
<td>emirpasic/gods</td>
<td>43,898,088</td>
<td>26.29 ns/op</td>
<td>48 B/op</td>
<td>0 allocs/op</td>
</tr>
<tr>
<td rowspan=2>AddAtIndex</td>
<td>architagr/golang_collections</td>
<td>363,227,274</td>
<td>3.300</td>
<td>0 B/op</td>
<td>0 allocs/op</td>
</tr>
<tr>
<td>emirpasic/gods</td>
<td>812,377</td>
<td>262970 ns/op</td>
<td>41 B/op</td>
<td>0 allocs/op</td>
</tr>

<tr>
<td rowspan=2>Find/IndexOf</td>
<td>architagr/golang_collections</td>
<td>99,529,309</td>
<td>11.59</td>
<td>0 B/op</td>
<td>0 allocs/op</td>
</tr>
<tr>
<td>emirpasic/gods <br/>
<br/>
issue found when we use pointer data as elements: https://github.com/emirpasic/gods/issues/269</td>
<td>415,523,194</td>
<td>2.888 ns/op</td>
<td>0 B/op</td>
<td>0 allocs/op</td>
</tr>
</tbody>
</table>
