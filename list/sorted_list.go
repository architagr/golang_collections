package list

type Less[T any] func(left, right T) bool

type sortedList[T any] struct {
	data   []T
	lessFn Less[T]
}

func InitSortedList[T any](lessFn Less[T]) ISortedItratorList[T] {
	return &sortedList[T]{
		data:   make([]T, 0, ArrayListCapacity),
		lessFn: lessFn,
	}
}

func (l *sortedList[T]) binarySearch(data T, start, end int) int {
	if start >= end {
		if l.lessFn(l.data[start], data) {
			return start + 1
		}
		return start
	}
	mid := (start + end) / 2
	midEle, _ := l.Get(mid)

	// if data.Equal(midEle) {
	// 	return mid
	// }

	if l.lessFn(midEle, data) {
		return l.binarySearch(data, mid+1, end)
	}
	return l.binarySearch(data, start, mid-1)
}

func (l *sortedList[T]) Add(data T) (resultIndex int) {

	if len(l.data) == 0 {
		l.data = append(l.data, data)
		resultIndex = len(l.data) - 1
		return
	}
	index := l.binarySearch(data, 0, len(l.data)-1)
	l.data = append(l.data[:index+1], l.data[index:]...)
	l.data[index] = data

	return
}

func (l *sortedList[T]) Remove(data T) (removedIndex int, err error) {
	removedIndex = -1
	index := l.Find(data)
	if index >= 0 {
		l.removeElement(index)
		return
	}
	err = ErrDataNotFoundError
	return
}

func (l *sortedList[T]) RemoveAtIndex(index int) (data T, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	l.removeElement(index)
	return
}

func (l *sortedList[T]) Count() int {
	return len(l.data)
}

func (l *sortedList[T]) Get(index int) (data T, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	return
}

func (l *sortedList[T]) Set(index int, data T) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	l.RemoveAtIndex(index)
	l.Add(data)
	return nil
}

func (l *sortedList[T]) Find(data T) (index int) {
	index = l.binarySearch(data, 0, len(l.data)-1)

	// if l.data[index].Equal(data) {
	// 	return index
	// }
	return -1
}

func (l *sortedList[T]) Filter(f Filterfunc[T]) []T {
	result := make([]T, 0, l.Count())
	for _, val := range l.data {
		if f(val) {
			result = append(result, val)
		}
	}
	return result
}

func (l *sortedList[T]) RemoveAll(f Filterfunc[T]) []T {
	removedData, result := make([]T, 0, l.Count()), make([]T, 0, l.Count())
	for _, val := range l.data {
		if f(val) {
			removedData = append(removedData, val)
		} else {
			result = append(result, val)
		}
	}
	l.data = result
	return removedData
}

func (l *sortedList[T]) DeepCopy() ([]T, error) {
	result := make([]T, 0, l.Count())
	// for _, val := range l.data {
	// 	if val != nil {
	// 		data, ok := (val.Copy()).(deepCopy)
	// 		if ok {
	// 			result = append(result, data)
	// 		}
	// 	}
	// }
	return result, nil
}

func (l *sortedList[T]) removeElement(index int) {
	l.data = append(l.data[:index], l.data[index+1:]...)
}

func (l *sortedList[T]) validateIndex(index int) error {
	if index < 0 || len(l.data)-1 < index {
		return ErrInvalidIndex
	}
	return nil
}
