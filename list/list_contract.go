package list

type IDeepCopy[T any] interface {
	*T
	deepCopyBaseInterface
}

type deepCopyBaseInterface interface {
	Copy() interface{}
	Equal(val interface{}) bool
}

type ISortedList[T any, deepCopy IDeepCopy[T]] interface {
	// Add adds data to the end of the list
	Add(data deepCopy) (resultIndex int)
	// Remove removes data from the list and returns error if the data is not found
	Remove(data deepCopy) (removedIndex int, err error)
	// RemoveAtIndex removes data at the index, if the index is not valid then returns error
	RemoveAtIndex(index int) (data deepCopy, err error)
	// Count return the count of elements in the list
	Count() int
	// Get gets data at the index, if index is not valid then it returns error
	Get(index int) (data deepCopy, err error)
	// Set updates the data at the index, if index is not valid then returns error
	Set(index int, data deepCopy) error
	// Find helps to get the first occourance if the data that matches the input and returns index
	// if index is -1 then the data is not found
	Find(data deepCopy) (index int)
}

type IList[T any, deepCopy IDeepCopy[T]] interface {
	ISortedList[T, deepCopy]
	// AddAtIndex adds data to the index in the list and shifts all data to right
	// if the index is out of bound then return error
	AddAtIndex(index int, data deepCopy) (err error)
}

type filterfunc[T any, deepCopy IDeepCopy[T]] func(data deepCopy) bool

type IItratorList[T any, deepCopy IDeepCopy[T]] interface {
	IList[T, deepCopy]
	// Filter helps to get the all data that matches according to the filter func and also returns index
	Filter(f filterfunc[T, deepCopy]) []deepCopy
	// DeepCopy this is used to create a copy of the list
	DeepCopy() []deepCopy
}

type ISortedItratorList[T any, deepCopy IDeepCopy[T]] interface {
	ISortedList[T, deepCopy]
	// Filter helps to get the all data that matches according to the filter func and also returns index
	Filter(f filterfunc[T, deepCopy]) []deepCopy
	// DeepCopy this is used to create a copy of the list
	DeepCopy() []deepCopy
}

type IIndexedItratorList[T any, deepCopy IDeepCopy[T]] interface {
	IItratorList[T, deepCopy]
	// FindByIndexedKey get the first occourance if the data that matches according to the filter func and also returns index,
	// this helps in fast search ad it will do a binary search on the indexKey
	// this is similar to Non-Clustered Index in database
	FindByIndexedKey(indexKey string, key string) (data T, index int)
	// FindByIndexedKey get the first occourance if the data that matches according to the filter func and also returns index,
	// this helps in fast search ad it will do a binary search on the indexKey
	// this is similar to Non-Clustered Index in database
	FilterByIndexedKey(indexKey string, key string) []T
}
