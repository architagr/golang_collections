package list

var (
	ArrayListCapacity int = 100
	growthPercentage      = float32(1.0)  // growth by 100%
	shrinkPercentage      = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

type arrayList[deepCopy IDeepCopy[T], T any] struct {
	data []deepCopy
	size int
}

func InitArrayList[deepCopy IDeepCopy[T], T any](data ...deepCopy) IItratorList[deepCopy, T] {
	size := 0
	if len(data) == 0 {
		data = make([]deepCopy, ArrayListCapacity)
	} else {
		size = len(data)
	}
	l := &arrayList[deepCopy, T]{
		data: data,
		size: size,
	}
	return l
}

func (l *arrayList[deepCopy, T]) Add(data deepCopy) (resultIndex int) {
	l.growArrayList(1)
	l.data[l.size] = data
	l.size++
	resultIndex = l.size - 1
	return
}

func (l *arrayList[deepCopy, T]) AddAtIndex(index int, data deepCopy) (err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	l.growArrayList(1)
	l.size++
	copy(l.data[index+1:], l.data[index:l.size-1])
	l.data[index] = data
	return
}

func (l *arrayList[deepCopy, T]) Remove(data deepCopy) (removedIndex int, err error) {
	removedIndex = -1
	for i := 0; i < l.Count(); i++ {
		if data.Equal(l.data[i]) {
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
	removedData := make([]deepCopy, 0, l.Count())
	for i := 0; i < l.Count(); {
		val := l.data[i]
		if f(val) {
			removedData = append(removedData, val)
			l.removeElement(i)
			continue
		}
		i++
	}
	return removedData
}

func (l *arrayList[deepCopy, T]) Count() int {
	return l.size
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
	for i := 0; i < l.size; i++ {
		if data.Equal(l.data[i]) {
			return i
		}
	}
	return -1
}

func (l *arrayList[deepCopy, T]) Filter(f Filterfunc[deepCopy, T]) []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for i := 0; i < l.size; i++ {
		if f(l.data[i]) {
			result = append(result, l.data[i])
		}
	}
	return result
}

func (l *arrayList[deepCopy, T]) DeepCopy() []deepCopy {
	result := make([]deepCopy, 0, l.Count())
	for i := 0; i < l.Count(); i++ {
		if l.data[i] != nil {
			data, ok := (l.data[i].Copy()).(deepCopy)
			if ok {
				result = append(result, data)
			}
		}
	}
	return result
}

func (l *arrayList[deepCopy, T]) removeElement(index int) {
	copy(l.data[index:], l.data[index+1:l.Count()])
	l.size--
	l.shrinkArrayList()
}

func (l *arrayList[deepCopy, T]) validateIndex(index int) error {
	if index < 0 || l.Count()-1 < index {
		return errInvalidIndex
	}
	return nil
}

func (l *arrayList[deepCopy, T]) resizeList(cap int) {
	newElements := make([]deepCopy, cap)
	copy(newElements, l.data)
	l.data = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (l *arrayList[deepCopy, T]) growArrayList(n int) {
	// When capacity is reached, grow by a percentage of growthPercentage and add number of elements
	currentCapacity := cap(l.data)
	if l.Count()+n >= currentCapacity {
		newCapacity := int((1 + growthPercentage) * float32(currentCapacity+n))
		l.resizeList(newCapacity)
	}
}

// Shrink the array if necessary,
// basically when size is shrinkPercentage % of current capacity
func (l *arrayList[deepCopy, T]) shrinkArrayList() {
	if shrinkPercentage == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(l.data)
	if l.Count() <= int(float32(currentCapacity)*shrinkPercentage) {
		l.resizeList(l.Count())
	}
}
