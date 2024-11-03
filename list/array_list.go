package list

var (
	ArrayListCapacity int = 100
)

type arrayList[deepCopy IDeepCopy[T], T any] struct {
	data []deepCopy
}

func InitArrayList[deepCopy IDeepCopy[T], T any](data ...deepCopy) IItratorList[deepCopy, T] {
	if len(data) == 0 {
		data = make([]deepCopy, 0, ArrayListCapacity)
	}
	l := &arrayList[deepCopy, T]{
		data: data,
	}
	return l
}

func (l *arrayList[deepCopy, T]) Add(data deepCopy) (resultIndex int) {
	l.data = append(l.data, data)
	resultIndex = len(l.data) - 1
	return
}

func (l *arrayList[deepCopy, T]) AddAtIndex(index int, data deepCopy) (err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	l.data = append(l.data[:index+1], l.data[index:]...)
	l.data[index] = data
	return
}

func (l *arrayList[deepCopy, T]) Remove(data deepCopy) (removedIndex int, err error) {
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

func (l *arrayList[deepCopy, T]) RemoveAtIndex(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	l.removeElement(index)
	return
}

func (l *arrayList[deepCopy, T]) RemoveAll(f Filterfunc[deepCopy, T]) []deepCopy {
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

func (l *arrayList[deepCopy, T]) Count() int {
	return len(l.data)
}

func (l *arrayList[deepCopy, T]) Get(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	return
}

func (l *arrayList[deepCopy, T]) Set(index int, data deepCopy) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	l.data[index] = data
	return nil
}

func (l *arrayList[deepCopy, T]) Find(data deepCopy) (index int) {
	index = -1
	for i, val := range l.data {
		if data.Equal(val) {
			index = i
			break
		}
	}
	return
}

func (l *arrayList[deepCopy, T]) Filter(f Filterfunc[deepCopy, T]) []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for _, val := range l.data {
		if f(val) {
			result = append(result, val)
		}
	}
	return result
}

func (l *arrayList[deepCopy, T]) DeepCopy() []deepCopy {
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

func (l *arrayList[deepCopy, T]) removeElement(index int) {
	l.data = append(l.data[:index], l.data[index+1:]...)
}

func (l *arrayList[deepCopy, T]) validateIndex(index int) error {
	if index < 0 || len(l.data)-1 < index {
		return errInvalidIndex
	}
	return nil
}
