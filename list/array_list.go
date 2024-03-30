package list

var (
	ArrayListBuffer   int = 10
	ArrayListCapacity int = 100
)

type arrayList[T any, deepCopy IDeepCopy[T]] struct {
	data []deepCopy
}

func InitArrayList[T any, deepCopy IDeepCopy[T]](data ...deepCopy) IItratorList[T, deepCopy] {
	l := &arrayList[T, deepCopy]{
		data: make([]deepCopy, 0, ArrayListCapacity),
	}
	if len(data) != 0 {
		l.data = data
	}
	return l
}

func (l *arrayList[T, deepCopy]) Add(data deepCopy) (resultIndex int) {
	defer l.adjustCapacity()
	l.data = append(l.data, data)
	resultIndex = len(l.data) - 1
	return
}

func (l *arrayList[T, deepCopy]) AddAtIndex(index int, data deepCopy) (err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	defer l.adjustCapacity()
	l.data = append(l.data[:index+1], l.data[index:]...)
	l.data[index] = data
	return
}

func (l *arrayList[T, deepCopy]) Remove(data deepCopy) (removedIndex int, err error) {
	removedIndex = -1
	for i, val := range l.data {
		if data.Equal(val) {
			l.removeElement(i)
			removedIndex = i
			return
		}
	}
	err = errDataNotFoundError
	return
}

func (l *arrayList[T, deepCopy]) RemoveAtIndex(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	l.removeElement(index)
	return
}

func (l *arrayList[T, deepCopy]) RemoveAll(f filterfunc[T, deepCopy]) []deepCopy {
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

func (l *arrayList[T, deepCopy]) Count() int {
	return len(l.data)
}

func (l *arrayList[T, deepCopy]) Get(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	return
}

func (l *arrayList[T, deepCopy]) Set(index int, data deepCopy) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	l.data[index] = data
	return nil
}

func (l *arrayList[T, deepCopy]) Find(data deepCopy) (index int) {
	index = -1
	for i, val := range l.data {
		if data.Equal(val) {
			index = i
			break
		}
	}
	return
}

func (l *arrayList[T, deepCopy]) Filter(f filterfunc[T, deepCopy]) []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for _, val := range l.data {
		if f(val) {
			result = append(result, val)
		}
	}
	return result
}

func (l *arrayList[T, deepCopy]) DeepCopy() []deepCopy {
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

func (l *arrayList[T, deepCopy]) removeElement(index int) {
	l.data = append(l.data[:index], l.data[index+1:]...)
}

func (l *arrayList[T, deepCopy]) adjustCapacity() {
	if cap(l.data)-len(l.data) <= ArrayListBuffer {
		newArray := make([]deepCopy, 0, cap(l.data)+ArrayListCapacity)
		l.data = append(newArray, l.data...)
	}
}

func (l *arrayList[T, deepCopy]) validateIndex(index int) error {
	if index < 0 || len(l.data)-1 < index {
		return errInvalidIndex
	}
	return nil
}
