package list

import (
	"bytes"
	"reflect"
)

// type IDeepCopy[T any] interface {
// 	*T
// 	deepCopyBaseInterface
// }
// type deepCopyBaseInterface interface {
// 	Copy() interface{}
// 	Equal(val interface{}) bool
// }

type IDeepCopy[T any] interface {
	Copy() interface{}
	Equal(val interface{}) bool
}

func AreEqual(data1, data2 any) bool {
	if data1 == nil || data2 == nil {
		return data1 == data2
	}

	exp, ok := data1.([]byte)
	if !ok {
		return reflect.DeepEqual(data1, data2)
	}

	act, ok := data2.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

func CheckImplementsDeepCopy[T any](data T) (IDeepCopy[T], bool) {
	dataKind := reflect.ValueOf(data).Kind()
	var o IDeepCopy[T]
	if dataKind != reflect.Pointer && reflect.TypeOf(data) != reflect.TypeOf(o) {
		return nil, false
	}
	v, ok := interface{}(&data).(IDeepCopy[T])
	return v, ok
}

type IBaseList[T any] interface {
	// Add adds data to the end of the list
	Add(data T) (resultIndex int)
	// Remove removes data from the list and returns error if the data is not found
	Remove(data T) (removedIndex int, err error)
	// RemoveAtIndex removes data at the index, if the index is not valid then returns error
	RemoveAtIndex(index int) (data T, err error)
	// Count return the count of elements in the list
	Count() int
	// Get gets data at the index, if index is not valid then it returns error
	Get(index int) (data T, err error)
	// Set updates the data at the index, if index is not valid then returns error
	Set(index int, data T) error
	// Find helps to get the first occourance if the data that matches the input and returns index
	// if index is -1 then the data is not found
	Find(data T) (index int)
	// RemoveAll function removes al elements for which the f returns true,
	// this also returns all elements removed
	RemoveAll(f Filterfunc[T]) []T
}

type IList[T any] interface {
	IBaseList[T]
	// AddAtIndex adds data to the index in the list and shifts all data to right
	// if the index is out of bound then return error
	AddAtIndex(index int, data T) (err error)
}

type Mapfunc[T1 any, T2 any] func(data T1) T2

type Filterfunc[T any] func(data T) bool

type IItratorList[T any] interface {
	IList[T]
	// Filter helps to get the all data that matches according to the filter func and also returns index
	Filter(f Filterfunc[T]) []T
	// DeepCopy this is used to create a copy of the list
	DeepCopy() ([]T, error)
}

type ISortedItratorList[T any] interface {
	IBaseList[T]
	// Filter helps to get the all data that matches according to the filter func and also returns index
	Filter(f Filterfunc[T]) []T
	// DeepCopy this is used to create a copy of the list
	DeepCopy() ([]T, error)
}

type IIndexedItratorList[T any] interface {
	IItratorList[T]
	// FindByIndexedKey get the first occourance if the data that matches according to the filter func and also returns index,
	// this helps in fast search ad it will do a binary search on the indexKey
	// this is similar to Non-Clustered Index in database
	FindByIndexedKey(indexKey string, key string) (data T, index int)
	// FindByIndexedKey get the first occourance if the data that matches according to the filter func and also returns index,
	// this helps in fast search ad it will do a binary search on the indexKey
	// this is similar to Non-Clustered Index in database
	FilterByIndexedKey(indexKey string, key string) []T
}
