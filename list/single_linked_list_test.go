package list

import (
	"testing"
)

type Integer int

func InitInteger(val int) *Integer {
	x := Integer(val)
	return &x
}
func (obj *Integer) Copy() interface{} {
	return obj
}

func (obj *Integer) Equal(val interface{}) bool {
	x := val.(*Integer)
	return int(*x) == int(*obj)
}

var singleLinkedListSampleData = []int{1, 2, 3}

func GetSingleLinkedList() IList[Integer, *Integer] {
	obj := InitSingleLinkedList[Integer]()
	// adding dummy data
	for _, val := range []int{1, 2, 3} {
		intVal := InitInteger(val)
		obj.Add(intVal)
	}
	return obj
}

func TestInitilizationSingleLinkedList(t *testing.T) {
	obj := InitSingleLinkedList[Integer]()

	if obj == nil || obj.Count() != 0 {
		t.Errorf("Single linked list not initilized correctly")
	}
}

// #region test Add node method
func TestAddNodeToNewSingleLinkedList(t *testing.T) {
	obj := InitSingleLinkedList[Integer]()
	val := InitInteger(1)
	obj.Add(val)
	if obj.Count() != 1 {
		t.Errorf("New Single linked list did not add 1 data correctly, expected element count is 1, got %d", obj.Count())
	}

	data, err := obj.Get(0)
	if err != nil {
		t.Fatalf("Single Linked List should not return error if trying to access valud index")
	}

	if int(*data) != 1 {
		t.Errorf("Single Linked List should return valid data")
	}
}

func TestAddNodeToSingleLinkedList(t *testing.T) {
	obj := GetSingleLinkedList()
	if obj.Count() != len(singleLinkedListSampleData) {
		t.Errorf("adding multiple node to single lisnked list failed, expected count %d, got:%d", len(singleLinkedListSampleData), obj.Count())
	}
	for validatingIndex, val := range singleLinkedListSampleData {
		data, err := obj.Get(validatingIndex)
		if err != nil {
			t.Fatalf("Single Linked List should not return error if trying to access valud index")
		}

		if int(*data) != val {
			t.Errorf("Single Linked List should return valid data")
		}
	}
}

// #endregion

// #region test AddAtIndex
func TestAddNodeAtInvalidIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	err := obj.AddAtIndex(len(singleLinkedListSampleData)+10, intVal)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to add at a invalid index")
	}
}

func TestAddNodeAtNegativeIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	err := obj.AddAtIndex(-1, intVal)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to add at a invalid index")
	}
}
func TestAddNodeAfterLastNode(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	err := obj.AddAtIndex(len(singleLinkedListSampleData), intVal)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to add at a invalid index")
	}
}

func TestAddNodeAtValidIndexInBetween(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	index := len(singleLinkedListSampleData) / 2
	err := obj.AddAtIndex(index, intVal)
	if err != nil {
		t.Errorf("Single Linked List should not return error if we try to add at a valid index")
	}
	data, _ := obj.Get(index)
	if int(*data) != 10 {
		t.Errorf("Single Linked List should return valid data for the index")
	}
}

func TestAddNodeAtValidIndexAtBeginning(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	index := 0
	err := obj.AddAtIndex(index, intVal)
	if err != nil {
		t.Errorf("Single Linked List should not return error if we try to add at a valid index")
	}
	data, _ := obj.Get(index)
	if int(*data) != 10 {
		t.Errorf("Single Linked List should return valid data for the index")
	}

}

func TestAddNodeAtValidIndexAtEnd(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	index := len(singleLinkedListSampleData) - 1
	err := obj.AddAtIndex(index, intVal)
	if err != nil {
		t.Errorf("Single Linked List should not return error if we try to add at a valid index")
	}
	data, _ := obj.Get(index)
	if int(*data) != 10 {
		t.Errorf("Single Linked List should return valid data for the index")
	}
	if obj.Count() != len(singleLinkedListSampleData)+1 {
		t.Errorf("Single Linked List error when adding at end")
	}
}

// #endregion

// #region test Get
func TestGetNodeWithValidIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	for _, validatingIndex := range []int{1, 0} {
		data, err := obj.Get(validatingIndex)
		if err != nil {
			t.Fatalf("Single Linked List should not return error if trying to access valid index")
		}

		if int(*data) != singleLinkedListSampleData[validatingIndex] {
			t.Errorf("Single Linked List should return valid data")
		}
	}
	if obj.Count() == 0 {
		t.Errorf("Single Linked List get function also removing nodes")
	}
}

func TestGetNodeWithNegativeIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	_, err := obj.Get(-1)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to access negitive index")
	}
}
func TestGetNodeWithOutofBoundIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	_, err := obj.Get(len(singleLinkedListSampleData) + 10)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to access index that does not exist")
	}
}

func TestGetNodeInEmptyList(t *testing.T) {
	obj := InitSingleLinkedList[Integer]()
	_, err := obj.Get(0)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to access any index when list is empty")
	}
}

// #endregion

// #region test Set
func TestSetNodeWithValidIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	err := obj.Set(0, intVal)
	if err != nil {
		t.Fatalf("Single Linked List should not return error if trying to set valid index")
	}
	data, err := obj.Get(0)
	if err != nil {
		t.Fatalf("Single Linked List should not return error if trying to access valid index")
	}

	if int(*data) != 10 {
		t.Errorf("Single Linked List should return valid data")
	}
}

func TestSetNodeWithNegativeIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	err := obj.Set(-1, intVal)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to access negitive index")
	}
}
func TestSetNodeWithOutofBoundIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	intVal := InitInteger(10)
	err := obj.Set(len(singleLinkedListSampleData)+10, intVal)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to set index that does not exist")
	}
}

func TestSetNodeInEmptyList(t *testing.T) {
	obj := InitSingleLinkedList[Integer]()
	intVal := InitInteger(10)
	err := obj.Set(0, intVal)
	if err == nil {
		t.Errorf("Single Linked List should return error if we try to set any index when list is empty")
	}
}

// #endregion

// #region test Remove
func TestRemovingValidData(t *testing.T) {
	obj := GetSingleLinkedList()
	removingIndex := 0
	data, _ := obj.Get(removingIndex)
	index, err := obj.Remove(data)
	if err != nil {
		t.Errorf("Single Linked List error while deleting a node")
	}
	if index != removingIndex {
		t.Errorf("Single Linked List removed wrong node")
	}
}
func TestRemovingDataThatDoesNotExist(t *testing.T) {
	obj := GetSingleLinkedList()
	val := InitInteger(10)
	_, err := obj.Remove(val)
	if err == nil {
		t.Errorf("Single Linked List should have error while deleting a node that does not exist")
	}
}

// #endregion

// #region test RemoveAtIndex
func TestRemoveValidIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	data, err := obj.RemoveAtIndex(1)
	if err != nil {
		t.Errorf("Single linked list should not throw error when trying to remove from valid index")
	}

	if int(*data) != singleLinkedListSampleData[1] {
		t.Errorf("Single linked list did not delete correct data")
	}
	if obj.Count() != len(singleLinkedListSampleData)-1 {
		t.Errorf("Single linked list did not delete data")
	}
}
func TestRemoveValidIndexFromBeginning(t *testing.T) {
	obj := GetSingleLinkedList()
	data, err := obj.RemoveAtIndex(0)
	if err != nil {
		t.Errorf("Single linked list should not throw error when trying to remove from valid index")
	}

	if int(*data) != singleLinkedListSampleData[0] {
		t.Errorf("Single linked list did not delete correct data")
	}

	if obj.Count() != len(singleLinkedListSampleData)-1 {
		t.Errorf("Single linked list did not delete data")
	}
}

func TestRemoveValidIndexFromEnd(t *testing.T) {
	obj := GetSingleLinkedList()
	data, err := obj.RemoveAtIndex(len(singleLinkedListSampleData) - 1)
	if err != nil {
		t.Errorf("Single linked list should not throw error when trying to remove from valid index")
	}

	if int(*data) != singleLinkedListSampleData[len(singleLinkedListSampleData)-1] {
		t.Errorf("Single linked list did not delete correct data")
	}

	if obj.Count() != len(singleLinkedListSampleData)-1 {
		t.Errorf("Single linked list did not delete data")
	}
}

func TestRemoveInValidNegitiveIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	_, err := obj.RemoveAtIndex(-1)
	if err == nil {
		t.Errorf("Single linked list should throw error when trying to remove from invalid index")
	}
}

func TestRemoveInValidIndex(t *testing.T) {
	obj := GetSingleLinkedList()
	_, err := obj.RemoveAtIndex(obj.Count() + 10)
	if err == nil {
		t.Errorf("Single linked list should throw error when trying to remove from invalid index")
	}
}

// #endregion

// #region test Find
func TestFindValidData(t *testing.T) {
	obj := GetSingleLinkedList()
	index := obj.Find(func(val *Integer) bool {
		x := InitInteger(2)
		return val.Equal(x)
	})

	if index != 1 {
		t.Errorf("Single linked list should be able to find a valid data")
	}

}
func TestFindInValidData(t *testing.T) {
	obj := GetSingleLinkedList()
	index := obj.Find(func(val *Integer) bool {
		x := InitInteger(20)
		return val.Equal(x)
	})

	if index != -1 {
		t.Errorf("Single linked list should return -1 if data is not found")
	}

}

// #endregion
