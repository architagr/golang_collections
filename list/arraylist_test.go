package list

import (
	"fmt"
	"testing"
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

// #region test initilization
// test if we are able to initilize a blank array list and have no element
func TestInitilizationOfBlankItratorArrayList(t *testing.T) {
	list := InitArrayList[*user]()

	if list.Count() != 0 {
		t.Fatalf("list not intilized")
	}
}

func TestInitilizationOfItratorArrayListWithValue(t *testing.T) {
	users := []*user{
		{
			id:   1,
			name: "test name 1",
		},
		{
			id:   2,
			name: "test name 2",
		},
	}
	list := InitArrayList(users...)

	if list.Count() != len(users) {
		t.Fatalf("list not intilized")
	}
}

// #endregion

// #region test Add
func TestAddingNewDataPositive(t *testing.T) {
	list := InitArrayList[*user]()
	input := &user{
		id:   1,
		name: "test 1",
	}
	list.Add(input)
	data, err := list.Get(0)
	if err != nil {
		t.Fatalf("failed to get new data")
	}
	fmt.Printf("received data %+v", data)
	if !data.Equal(input) {
		t.Fatalf("expect data to be same as the value inserted")
	}
}
func TestAddingNewDataPositiveMoreData(t *testing.T) {
	ArrayListCapacity = 10
	list := InitArrayList[*user]()
	inputs := getSampleData()
	for _, val := range inputs {
		list.Add(val)
	}
	data, err := list.Get(0)
	if err != nil {
		t.Fatalf("failed to get new data")
	}
	fmt.Printf("received data %+v", data)
	if !data.Equal(inputs[0]) {
		t.Fatalf("expect data to be same as the value inserted")
	}
}

// #endregion

// #region test AddAtIndex
func getSampleData() []*user {
	return []*user{
		{
			id:   1,
			name: "test 1",
		},
		{
			id:   2,
			name: "test 2",
		},
		{
			id:   3,
			name: "test 3",
		},
		{
			id:   4,
			name: "test 4",
		},
		{
			id:   5,
			name: "test 5",
		},
	}
}
func TestAddingAtIndexNewDataInvalidIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	input := &user{
		id:   10,
		name: "test 10",
	}
	err := list.AddAtIndex(10, input)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wrong index")
	}

}
func TestAddingAtIndexNewDataNegativeIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	input := &user{
		id:   10,
		name: "test 10",
	}
	err := list.AddAtIndex(-1, input)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wring index")
	}

}

func TestAddingAtIndexNewData(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	input := &user{
		id:   10,
		name: "test 10",
	}
	err := list.AddAtIndex(2, input)
	if err != nil {
		t.Fatalf("Array list should be able to add value at a index")
	}

	data, _ := list.Get(2)
	if !data.Equal(input) {
		t.Fatalf("expect data to be same as the value inserted")
	}

}

// #endregion
// #region test get
func TestGetWrongIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	_, err := list.Get(100)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wrong index")
	}
}
func TestGetNegativeIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	_, err := list.Get(-1)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wrong index")
	}
}

func TestGet(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	data, err := list.Get(1)
	if err != nil {
		t.Fatalf("Array list return data")
	}

	if !data.Equal(getSampleData()[1]) {
		t.Fatalf("Array list is not returning correct value")
	}
}

// #endregion
// #region test remove at index
func TestRemoveAtIndexWrongIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	_, err := list.RemoveAtIndex(100)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}
func TestRemoveAtIndexNegativeIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	_, err := list.RemoveAtIndex(-1)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}

func TestRemoveAtIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	data, err := list.RemoveAtIndex(1)
	if err != nil {
		t.Fatalf("Array list return data")
	}

	if !data.Equal(getSampleData()[1]) {
		t.Fatalf("Array list is not returning correct value")
	}
	if list.Count()+1 != len(getSampleData()) {
		t.Fatalf("data was not removed")
	}
}

// #endregion

// #region test remove

func TestRemoveNegative(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	_, err := list.Remove(&user{
		id:   100,
		name: "test",
	})
	if err != errDataNotFoundError {
		t.Fatalf("Array list should through error if removing wrong index")
	}
}

func TestRemove(t *testing.T) {
	sampleData := getSampleData()
	list := InitArrayList[*user](sampleData...)
	index, err := list.Remove(sampleData[1])
	if err != nil {
		t.Fatalf("Array list remove is not removing a valid data")
	}

	if index != 1 {
		t.Fatalf("Array list is not returning correct index that is removed")
	}
	if list.Count()+1 != len(getSampleData()) {
		t.Fatalf("data was not removed")
	}
}

// #endregion

// #region test set
func TestSetWrongIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	err := list.Set(100, &user{
		id:   100,
		name: "test",
	})
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}
func TestSetNegativeIndex(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	err := list.Set(-1, &user{
		id:   100,
		name: "test",
	})
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}

func TestSet(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	usr := &user{
		id:   100,
		name: "test",
	}
	err := list.Set(1, usr)
	if err != nil {
		t.Fatalf("Array list return data")
	}
	newData, _ := list.Get(1)

	if !newData.Equal(usr) {
		t.Fatalf("Array list has not set correct value")
	}
}

// #endregion
// #region test find
func TestFindNegative(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	index := list.Find(&user{
		id:   1,
		name: "test",
	})
	if index != -1 {
		t.Fatalf("Array list should return -1 as a response of find")
	}
}

func TestFind(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	index := list.Find(&user{
		id:   1,
		name: "test 1",
	})
	if index == -1 {
		t.Fatalf("Array list should not return -1 as a response of find")
	}
}

// #endregion

// #region test filter
func TestFilterNegative(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	response := list.Filter(func(val *user) bool {
		return val.id == 1 && val.name == "test"
	})
	if len(response) != 0 {
		t.Fatalf("Array list should return no data as a response of filter")
	}
}

func TestFilter(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	response := list.Filter(func(val *user) bool {
		return val.id == 1
	})
	if len(response) != 1 {
		t.Fatalf("Array list should only return 1 data")
	}
}

// #endregion

func TestDeepCopy(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	response := list.DeepCopy()
	if len(response) != len(getSampleData()) {
		t.Fatalf("Array list return same data")
	}
}

// #region test RemoveAll

func TestRemoveAll(t *testing.T) {
	list := InitArrayList[*user](getSampleData()...)
	removedData := list.RemoveAll(func(val *user) bool {
		return val.id%2 == 0
	})
	for _, val := range removedData {
		if val.id%2 != 0 {
			t.Fatalf("removed some data that was not supposed to be removed")
		}
	}
	response := list.DeepCopy()
	for _, val := range response {
		if val.id%2 == 0 {
			t.Fatalf("did not remove data that was not supposed to be removed")
		}
	}
}

// #endregion
