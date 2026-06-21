package list

import "fmt"

var (
	ArrayListCapacity int = 100
	growthPercentage      = float32(1.0)  // growth by 100%
	shrinkPercentage      = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

type arrayList[T any] struct {
	data               []T
	deepCopyData       []IDeepCopy[T]
	size               int
	implementsDeepCopy *bool
}

func InitArrayList[T any](data ...T) IItratorList[T] {
	size := 0

	if len(data) == 0 {
		data = make([]T, ArrayListCapacity)
	} else {
		size = len(data)

	}
	l := &arrayList[T]{
		data: data,
		size: size,
	}
	l.setImplementsDeepCopy()
	return l
}

func (l *arrayList[T]) setImplementsDeepCopy() {
	var o T
	var implementsDeepCopy *bool
	obj, ok := CheckImplementsDeepCopy(o)
	obj, ok = CheckImplementsDeepCopy(l.data[0])
	fmt.Println(obj)
	implementsDeepCopy = &ok
	l.implementsDeepCopy = implementsDeepCopy
}

func (l *arrayList[T]) Add(data T) (resultIndex int) {
	l.growArrayList(1)
	l.data[l.size] = data
	l.size++
	resultIndex = l.size - 1
	return
}

func (l *arrayList[T]) AddAtIndex(index int, data T) (err error) {
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

func (l *arrayList[T]) Remove(data T) (removedIndex int, err error) {
	removedIndex = -1
	for i := 0; i < l.Count(); i++ {
		if AreEqual(data, l.data[i]) {
			l.removeElement(i)
			removedIndex = i
			return
		}
	}
	err = ErrDataNotFoundError
	return
}

func (l *arrayList[T]) RemoveAtIndex(index int) (data T, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	l.removeElement(index)
	return
}

func (l *arrayList[T]) RemoveAll(f Filterfunc[T]) []T {
	removedData := make([]T, 0, l.Count())
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

func (l *arrayList[T]) Count() int {
	return l.size
}

func (l *arrayList[T]) Get(index int) (data T, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	data = l.data[index]
	return
}

func (l *arrayList[T]) Set(index int, data T) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	l.data[index] = data
	return nil
}

func (l *arrayList[T]) Find(data T) (index int) {
	v, implementsDeepCopy := CheckImplementsDeepCopy(data)
	index = -1
	for i := 0; i < l.size; i++ {
		if implementsDeepCopy {
			if v.Equal(l.data[i]) {
				index = i
				break
			}
		} else {
			if AreEqual(data, l.data[i]) {
				index = i
				break
			}
		}
	}
	return index
}

func (l *arrayList[T]) Filter(f Filterfunc[T]) []T {
	result := make([]T, 0, l.Count())
	for i := 0; i < l.size; i++ {
		if f(l.data[i]) {
			result = append(result, l.data[i])
		}
	}
	return result
}

func (l *arrayList[T]) DeepCopy() ([]T, error) {
	if l.Count() == 0 {
		return nil, nil
	}
	if _, ok := CheckImplementsDeepCopy(l.data[0]); !ok {
		var o IDeepCopy[T]
		return nil, fmt.Errorf("%w %T", ErrDoesNotImplement, o)
	}
	result := make([]T, 0, l.Count())
	for i := 0; i < l.Count(); i++ {
		v, _ := CheckImplementsDeepCopy(l.data[i])
		if v != nil {
			data := v.Copy()
			o, _ := data.(T)
			result = append(result, o)
		}
	}
	return result, nil
}

func (l *arrayList[T]) removeElement(index int) {
	copy(l.data[index:], l.data[index+1:l.Count()])
	l.size--
	l.shrinkArrayList()
}

func (l *arrayList[T]) validateIndex(index int) error {
	if index < 0 || l.Count()-1 < index {
		return ErrInvalidIndex
	}
	return nil
}

func (l *arrayList[T]) resizeList(cap int) {
	newElements := make([]T, cap)
	copy(newElements, l.data)
	l.data = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (l *arrayList[T]) growArrayList(n int) {
	// When capacity is reached, grow by a percentage of growthPercentage and add number of elements
	currentCapacity := cap(l.data)
	if l.Count()+n >= currentCapacity {
		newCapacity := int((1 + growthPercentage) * float32(currentCapacity+n))
		l.resizeList(newCapacity)
	}
}

// Shrink the array if necessary,
// basically when size is shrinkPercentage % of current capacity
func (l *arrayList[T]) shrinkArrayList() {
	if shrinkPercentage == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(l.data)
	if l.Count() <= int(float32(currentCapacity)*shrinkPercentage) {
		l.resizeList(l.Count())
	}
}
