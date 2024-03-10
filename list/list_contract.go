package list

type IDeepCopy interface {
	Copy() interface{}
	Equal(val interface{}) bool
}

type IList[T IDeepCopy] interface {
	// Add adds data to the end of the list
	Add(data *T) (resultIndex int)
	// AddAtIndex adds data to the index in the list and shifts all data to right
	// if the index is out of bound then return error
	AddAtIndex(index int, data *T) (err error)
	// Remove removes data from the list and returns error if the data is not found
	Remove(data *T) (removedIndex int, err error)
	// RemoveAtIndex removes data at the index, if the index is not valid then returns error
	RemoveAtIndex(index int) (data *T, err error)
	// Count return the count of elements in the list
	Count() int
	// Get gets data at the index, if index is not valid then it returns error
	Get(index int) (data *T, err error)
	// Set updates the data at the index, if index is not valid then returns error
	Set(index int, data *T) error
	// Find helps to get the first occourance if the data that matches according to the filter func and also returns index
	// if index is -1 then the data is not found
	Find(f filterfunc[T]) (index int)
}

type filterfunc[T IDeepCopy] func(data *T) bool

type IItratorList[T IDeepCopy] interface {
	IList[T]
	// Filter helps to get the all data that matches according to the filter func and also returns index
	Filter(f filterfunc[T]) []T
	// DeepCopy this is used to create a copy of the list
	DeepCopy() []T
}

type IIndexedItratorList[T IDeepCopy] interface {
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
