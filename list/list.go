package list

type IList[T any] interface {
	Add(data T)
	AddRange(data ...T)
	Find(func(ele T, index int) bool) T
	Filter(func(ele T, index int) bool) IList[T]
	Some(func(ele T, index int) bool) bool
}

type List[T any] struct {
	data []T
}

func NewList[T any](data ...T) IList[T] {
	l := &List[T]{
		data: make([]T, 0, 100),
	}
	if len(data) != 0 {
		l.data = data
	}
	return l
}

func (l *List[T]) Add(data T) {
	l.data = append(l.data, data)
}
func (l *List[T]) AddRange(data ...T) {
	l.data = append(l.data, data...)
}
func (l *List[T]) Find(func(ele T, index int) bool) T {
	return l.data[0]
}
func (l *List[T]) Filter(func(ele T, index int) bool) IList[T] {
	return l
}
func (l *List[T]) Some(func(ele T, index int) bool) bool {
	return true
}
