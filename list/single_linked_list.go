package list

func InitSingleLinkedList[deepCopy IDeepCopy[T], T any]() IList[deepCopy, T] {
	return &singleLinkedList[deepCopy, T]{
		head:     nil,
		tail:     nil,
		indexMap: make(map[int]*singleLinkedListNode[deepCopy, T]),
	}
}

type singleLinkedListNode[deepCopy IDeepCopy[T], T any] struct {
	data deepCopy
	next *singleLinkedListNode[deepCopy, T]
}

func initSingleLinkedListNode[deepCopy IDeepCopy[T], T any](data deepCopy) *singleLinkedListNode[deepCopy, T] {
	return &singleLinkedListNode[deepCopy, T]{
		data: data,
		next: nil,
	}
}

type singleLinkedList[deepCopy IDeepCopy[T], T any] struct {
	head, tail *singleLinkedListNode[deepCopy, T]
	indexMap   map[int]*singleLinkedListNode[deepCopy, T]
}

func (l *singleLinkedList[deepCopy, T]) Add(data deepCopy) (resultIndex int) {
	newNode := initSingleLinkedListNode(data)
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.indexMap[l.Count()] = newNode
	return l.Count()
}

func (l *singleLinkedList[deepCopy, T]) AddAtIndex(index int, data deepCopy) (err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}

	newNode := initSingleLinkedListNode(data)

	for i := l.Count() - 1; i >= index; i-- {
		l.indexMap[i+1] = l.indexMap[i]
	}
	l.indexMap[index] = newNode

	if index == 0 {
		newNode.next = l.head
		l.head = newNode
	} else {
		newNode.next = l.indexMap[index+1]
		l.indexMap[index-1].next = newNode
	}

	return
}

func (l *singleLinkedList[deepCopy, T]) Remove(data deepCopy) (removedIndex int, err error) {
	temp := l.head
	removedIndex = 0
	for temp != nil {
		if data.Equal(temp.data) {
			_, err = l.RemoveAtIndex(removedIndex)
			return
		}
		removedIndex++
		temp = temp.next
	}
	removedIndex = -1
	err = errDataNotFoundError
	return
}
func (l *singleLinkedList[deepCopy, T]) RemoveAtIndex(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}

	if index == 0 {
		l.head = l.head.next
	} else if index+1 == l.Count() {
		l.indexMap[index-1].next = nil
		l.tail = l.indexMap[index-1]
	} else {
		l.indexMap[index-1].next = l.indexMap[index+1]
	}

	initialCount := l.Count()
	data = l.indexMap[index].data
	delete(l.indexMap, index)
	for i := index; i < initialCount-1; i++ {
		l.indexMap[i] = l.indexMap[i+1]
	}
	delete(l.indexMap, initialCount-1)
	return
}
func (l *singleLinkedList[deepCopy, T]) Count() int {
	return len(l.indexMap)
}
func (l *singleLinkedList[deepCopy, T]) validateIndex(index int) error {
	if index < 0 || index >= l.Count() {
		return errInvalidIndex
	}
	return nil
}
func (l *singleLinkedList[deepCopy, T]) Get(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	return l.indexMap[index].data, nil
}
func (l *singleLinkedList[deepCopy, T]) Set(index int, data deepCopy) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	node := l.indexMap[index]
	node.data = data
	return nil
}

func (l *singleLinkedList[deepCopy, T]) Find(data deepCopy) (index int) {
	temp := l.head
	index = 0
	for temp != nil {
		if temp.data.Equal(data) {
			break
		}
		index++
		temp = temp.next
	}
	if index == l.Count() {
		index = -1
	}
	return
}

func (l *singleLinkedList[deepCopy, T]) RemoveAll(f Filterfunc[deepCopy, T]) []deepCopy {
	removedData := make([]deepCopy, 0, l.Count())
	temp := l.head
	index := 0
	for temp != nil {
		d := temp.data
		temp = temp.next
		if f(d) {
			x, _ := l.RemoveAtIndex(index)
			removedData = append(removedData, x)
		} else {
			index++
		}

	}
	return removedData
}
