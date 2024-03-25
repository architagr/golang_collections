package list

func InitSingleLinkedList[T any, deepCopy IDeepCopy[T]]() IList[T, deepCopy] {
	return &singleLinkedList[T, deepCopy]{
		head:     nil,
		tail:     nil,
		indexMap: make(map[int]*singleLinkedListNode[T, deepCopy]),
	}
}

type singleLinkedListNode[T any, deepCopy IDeepCopy[T]] struct {
	data deepCopy
	next *singleLinkedListNode[T, deepCopy]
}

func initSingleLinkedListNode[T any, deepCopy IDeepCopy[T]](data deepCopy) *singleLinkedListNode[T, deepCopy] {
	return &singleLinkedListNode[T, deepCopy]{
		data: data,
		next: nil,
	}
}

type singleLinkedList[T any, deepCopy IDeepCopy[T]] struct {
	head, tail *singleLinkedListNode[T, deepCopy]
	indexMap   map[int]*singleLinkedListNode[T, deepCopy]
}

func (l *singleLinkedList[T, deepCopy]) Add(data deepCopy) (resultIndex int) {
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

func (l *singleLinkedList[T, deepCopy]) AddAtIndex(index int, data deepCopy) (err error) {
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

func (l *singleLinkedList[T, deepCopy]) Remove(data deepCopy) (removedIndex int, err error) {
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
func (l *singleLinkedList[T, deepCopy]) RemoveAtIndex(index int) (data deepCopy, err error) {
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
func (l *singleLinkedList[T, deepCopy]) Count() int {
	return len(l.indexMap)
}
func (l *singleLinkedList[T, deepCopy]) validateIndex(index int) error {
	if index < 0 || index >= l.Count() {
		return errInvalidIndex
	}
	return nil
}
func (l *singleLinkedList[T, deepCopy]) Get(index int) (data deepCopy, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	return l.indexMap[index].data, nil
}
func (l *singleLinkedList[T, deepCopy]) Set(index int, data deepCopy) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	node := l.indexMap[index]
	node.data = data
	return nil
}

func (l *singleLinkedList[T, deepCopy]) Find(data deepCopy) (index int) {
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
