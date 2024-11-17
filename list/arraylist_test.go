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
	userList := InitArrayList[*user]()

	if userList.Count() != 0 {
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
	userList := InitArrayList(users...)

	if userList.Count() != len(users) {
		t.Fatalf("list not intilized")
	}
}

// #endregion

// #region test Add
func TestAddingNewDataPositive(t *testing.T) {
	userList := InitArrayList[*user]()
	input := &user{
		id:   1,
		name: "test 1",
	}
	userList.Add(input)
	data, err := userList.Get(0)
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
	userList := InitArrayList(getSampleData()...)
	input := &user{
		id:   10,
		name: "test 10",
	}
	err := userList.AddAtIndex(10, input)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wrong index")
	}

}
func TestAddingAtIndexNewDataNegativeIndex(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	input := &user{
		id:   10,
		name: "test 10",
	}
	err := userList.AddAtIndex(-1, input)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wring index")
	}

}

func TestAddingAtIndexNewData(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	input := &user{
		id:   10,
		name: "test 10",
	}
	err := userList.AddAtIndex(2, input)
	if err != nil {
		t.Fatalf("Array list should be able to add value at a index")
	}

	data, _ := userList.Get(2)
	if !data.Equal(input) {
		t.Fatalf("expect data to be same as the value inserted")
	}

}

// #endregion
// #region test get
func TestGetWrongIndex(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	_, err := userList.Get(100)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wrong index")
	}
}
func TestGetNegativeIndex(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	_, err := userList.Get(-1)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if adding at wrong index")
	}
}

func TestGet(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	data, err := userList.Get(1)
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
	userList := InitArrayList(getSampleData()...)
	_, err := userList.RemoveAtIndex(100)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}
func TestRemoveAtIndexNegativeIndex(t *testing.T) {
	list := InitArrayList(getSampleData()...)
	_, err := list.RemoveAtIndex(-1)
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}

func TestRemoveAtIndex(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	data, err := userList.RemoveAtIndex(1)
	if err != nil {
		t.Fatalf("Array list return data")
	}

	if !data.Equal(getSampleData()[1]) {
		t.Fatalf("Array list is not returning correct value")
	}
	if userList.Count()+1 != len(getSampleData()) {
		t.Fatalf("data was not removed")
	}
}

// #endregion

// #region test remove
func TestRemoveNegative(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	_, err := userList.Remove(&user{
		id:   100,
		name: "test",
	})
	if err != errDataNotFoundError {
		t.Fatalf("Array list should through error if removing wrong index")
	}
}

func TestRemove(t *testing.T) {
	sampleData := getSampleData()
	userList := InitArrayList(sampleData...)
	index, err := userList.Remove(sampleData[1])
	if err != nil {
		t.Fatalf("Array list remove is not removing a valid data")
	}

	if index != 1 {
		t.Fatalf("Array list is not returning correct index that is removed")
	}
	if userList.Count()+1 != len(getSampleData()) {
		t.Fatalf("data was not removed")
	}
}

// #endregion

// #region test set
func TestSetWrongIndex(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	err := userList.Set(100, &user{
		id:   100,
		name: "test",
	})
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}
func TestSetNegativeIndex(t *testing.T) {
	list := InitArrayList(getSampleData()...)
	err := list.Set(-1, &user{
		id:   100,
		name: "test",
	})
	if err != errInvalidIndex {
		t.Fatalf("Array list should through error if removing at wrong index")
	}
}

func TestSet(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	usr := &user{
		id:   100,
		name: "test",
	}
	err := userList.Set(1, usr)
	if err != nil {
		t.Fatalf("Array list return data")
	}
	newData, _ := userList.Get(1)

	if !newData.Equal(usr) {
		t.Fatalf("Array list has not set correct value")
	}
}

// #endregion
// #region test find
func TestFindNegative(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	index := userList.Find(&user{
		id:   1,
		name: "test",
	})
	if index != -1 {
		t.Fatalf("Array list should return -1 as a response of find")
	}
}

func TestFind(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	index := userList.Find(&user{
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
	userList := InitArrayList(getSampleData()...)
	response := userList.Filter(func(val *user) bool {
		return val.id == 1 && val.name == "test"
	})
	if len(response) != 0 {
		t.Fatalf("Array list should return no data as a response of filter")
	}
}

func TestFilter(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	response := userList.Filter(func(val *user) bool {
		return val.id == 1
	})
	if len(response) != 1 {
		t.Fatalf("Array list should only return 1 data")
	}
}

// #endregion

func TestDeepCopy(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	response := userList.DeepCopy()
	if len(response) != len(getSampleData()) {
		t.Fatalf("Array list return same data")
	}
}

// #region test RemoveAll

func TestRemoveAll(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	removedData := userList.RemoveAll(func(val *user) bool {
		return val.id%2 == 0
	})
	for _, val := range removedData {
		if val.id%2 != 0 {
			t.Fatalf("removed some data that was not supposed to be removed")
		}
	}
	response := userList.DeepCopy()
	for _, val := range response {
		if val.id%2 == 0 {
			t.Fatalf("did not remove data that was not supposed to be removed")
		}
	}
}

func TestRemoveAll_Clear(t *testing.T) {
	userList := InitArrayList(getSampleData()...)
	removedData := userList.RemoveAll(func(val *user) bool {
		return val.id != 0
	})
	if userList.Count() != 0 || len(removedData) != len(getSampleData()) {
		t.Fatalf("did not remove data that was not supposed to be removed")
	}
}

func TestRemoveAll_Clear_NoShrink(t *testing.T) {
	shrinkPercentage = float32(0.0)
	userList := InitArrayList(getSampleData()...)
	removedData := userList.RemoveAll(func(val *user) bool {
		return val.id != 0
	})
	if userList.Count() != 0 || len(removedData) != len(getSampleData()) {
		t.Fatalf("did not remove data that was not supposed to be removed")
	}
}

// #endregion
