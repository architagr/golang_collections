package list

type Less[T any, deepCopy IDeepCopy[T]] func(left, right deepCopy) bool

type sortedList[T any, deepCopy IDeepCopy[T]] struct {
	data   []deepCopy
	lessFn Less[T, deepCopy]
}

func InitSortedList[T any, deepCopy IDeepCopy[T]](lessFn Less[T, deepCopy]) ISortedItratorList[T, deepCopy] {
	return &sortedList[T, deepCopy]{
		data:   make([]deepCopy, 0, ArrayListCapacity),
		lessFn: lessFn,
	}
}
func (l *sortedList[T, deepCopy]) binarySearch(data deepCopy, start, end int) int {
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

func (l *sortedList[T, deepCopy]) Add(data deepCopy) (resultIndex int) {
	defer l.adjustCapacity()

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

func (l *sortedList[T, deepCopy]) Remove(data deepCopy) (removedIndex int, err error) {
	removedIndex = -1
	index := l.Find(data)
	if index >= 0 {
		l.removeElement(index)
		return
	}
	err = errDataNotFoundError
	return
}

func (l *sortedList[T, deepCopy]) RemoveAtIndex(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	l.removeElement(index)
	return
}
func (l *sortedList[T, deepCopy]) Count() int {
	return len(l.data)
}
func (l *sortedList[T, deepCopy]) Get(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	return
}
func (l *sortedList[T, deepCopy]) Set(index int, data deepCopy) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	l.RemoveAtIndex(index)
	l.Add(data)
	return nil
}

func (l *sortedList[T, deepCopy]) Find(data deepCopy) (index int) {
	index = l.binarySearch(data, 0, len(l.data)-1)

	if l.data[index].Equal(data) {
		return index
	}
	return -1
}

func (l *sortedList[T, deepCopy]) Filter(f filterfunc[T, deepCopy]) []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for _, val := range l.data {
		if f(val) {
			result = append(result, val)
		}
	}
	return result
}
func (l *sortedList[T, deepCopy]) DeepCopy() []deepCopy {
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

func (l *sortedList[T, deepCopy]) removeElement(index int) {
	l.data = append(l.data[:index], l.data[index+1:]...)
}

func (l *sortedList[T, deepCopy]) adjustCapacity() {
	if cap(l.data)-len(l.data) <= ArrayListBuffer {
		newArray := make([]deepCopy, 0, cap(l.data)+ArrayListCapacity)
		l.data = append(newArray, l.data...)
	}
}
func (l *sortedList[T, deepCopy]) validateIndex(index int) error {
	if index < 0 || len(l.data)-1 < index {
		return errInvalidIndex
	}
	return nil
}
