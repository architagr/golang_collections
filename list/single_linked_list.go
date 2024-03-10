package list

import "fmt"

func InitSingleLinkedList[T IDeepCopy]() IList[T] {
	return &singleLinkedList[T]{
		head:     nil,
		tail:     nil,
		indexMap: make(map[int]*singleLinkedListNode[T]),
	}
}

type singleLinkedListNode[T IDeepCopy] struct {
	data *T
	next *singleLinkedListNode[T]
}

func initSingleLinkedListNode[T IDeepCopy](data *T) *singleLinkedListNode[T] {
	return &singleLinkedListNode[T]{
		data: data,
		next: nil,
	}
}

type singleLinkedList[T IDeepCopy] struct {
	head, tail *singleLinkedListNode[T]
	indexMap   map[int]*singleLinkedListNode[T]
}

func (l *singleLinkedList[T]) Add(data *T) (resultIndex int) {
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

func (l *singleLinkedList[T]) AddAtIndex(index int, data *T) (err error) {
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

func (l *singleLinkedList[T]) Remove(data *T) (removedIndex int, err error) {
	temp := l.head
	removedIndex = 0
	for temp != nil {
		if (*data).Equal(*temp.data) {
			_, err = l.RemoveAtIndex(removedIndex)
			break
		}
		removedIndex++
		temp = temp.next
	}
	if removedIndex == l.Count() {
		err = fmt.Errorf("data not found")
		removedIndex = -1
	}
	return
}
func (l *singleLinkedList[T]) RemoveAtIndex(index int) (data *T, err error) {
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
func (l *singleLinkedList[T]) Count() int {
	return len(l.indexMap)
}
func (l *singleLinkedList[T]) validateIndex(index int) error {
	if index < 0 || index >= l.Count() {
		return fmt.Errorf("invalid index")
	}
	return nil
}
func (l *singleLinkedList[T]) Get(index int) (data *T, err error) {
	err = l.validateIndex(index)
	if err != nil {
		return
	}
	return l.indexMap[index].data, nil
}
func (l *singleLinkedList[T]) Set(index int, data *T) error {
	err := l.validateIndex(index)
	if err != nil {
		return err
	}
	node := l.indexMap[index]
	node.data = data
	return nil
}

func (l *singleLinkedList[T]) Find(f filterfunc[T]) (index int) {
	temp := l.head
	index = 0
	for temp != nil {
		if f(temp.data) {
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
