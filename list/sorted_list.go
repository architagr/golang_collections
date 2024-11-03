package list

type Less[deepCopy IDeepCopy[T], T any] func(left, right deepCopy) bool

type sortedList[deepCopy IDeepCopy[T], T any] struct {
	data   []deepCopy
	lessFn Less[deepCopy, T]
}

func InitSortedList[deepCopy IDeepCopy[T], T any](lessFn Less[deepCopy, T]) ISortedItratorList[deepCopy, T] {
	return &sortedList[deepCopy, T]{
		data:   make([]deepCopy, 0, ArrayListCapacity),
		lessFn: lessFn,
	}
}

func (l *sortedList[deepCopy, T]) binarySearch(data deepCopy, start, end int) int {
	if start >= end {
		if l.lessFn(l.data[start], data) {
			return start + 1
		}
		return start
	}
	mid := (start + end) / 2
	midEle, _ := l.Get(mid)

	if data.Equal(midEle) {
		return mid
	}

	if l.lessFn(midEle, data) {
		return l.binarySearch(data, mid+1, end)
	}
	return l.binarySearch(data, start, mid-1)
}

func (l *sortedList[deepCopy, T]) Add(data deepCopy) (resultIndex int) {

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

func (l *sortedList[deepCopy, T]) Remove(data deepCopy) (removedIndex int, err error) {
	removedIndex = -1
	index := l.Find(data)
	if index >= 0 {
		l.removeElement(index)
		return
	}
	err = errDataNotFoundError
	return
}

func (l *sortedList[deepCopy, T]) RemoveAtIndex(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	l.removeElement(index)
	return
}

func (l *sortedList[deepCopy, T]) Count() int {
	return len(l.data)
}

func (l *sortedList[deepCopy, T]) Get(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	return
}

func (l *sortedList[deepCopy, T]) Set(index int, data deepCopy) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	l.RemoveAtIndex(index)
	l.Add(data)
	return nil
}

func (l *sortedList[deepCopy, T]) Find(data deepCopy) (index int) {
	index = l.binarySearch(data, 0, len(l.data)-1)

	if l.data[index].Equal(data) {
		return index
	}
	return -1
}

func (l *sortedList[deepCopy, T]) Filter(f Filterfunc[deepCopy, T]) []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for _, val := range l.data {
		if f(val) {
			result = append(result, val)
		}
	}
	return result
}

func (l *sortedList[deepCopy, T]) RemoveAll(f Filterfunc[deepCopy, T]) []deepCopy {
	removedData, result := make([]deepCopy, 0, l.Count()), make([]deepCopy, 0, l.Count())
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

func (l *sortedList[deepCopy, T]) DeepCopy() []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for _, val := range l.data {
		if val != nil {
			data, ok := (val.Copy()).(deepCopy)
			if ok {
				result = append(result, data)
			}
		}
	}
	return result
}

func (l *sortedList[deepCopy, T]) removeElement(index int) {
	l.data = append(l.data[:index], l.data[index+1:]...)
}

func (l *sortedList[deepCopy, T]) validateIndex(index int) error {
	if index < 0 || len(l.data)-1 < index {
		return errInvalidIndex
	}
	return nil
}
