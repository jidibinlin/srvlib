package doublelinked

import (
	"errors"
)

type LinkedList struct {
	header *LinkedNode
	tail   *LinkedNode
	len    uint32
}

func (l *LinkedList) Walk(fn func(node *LinkedNode) (bool, error)) error {
	node := l.header

	for node != nil {
		ok, err := fn(node)
		if err != nil {
			return err
		}

		if !ok {
			break
		}

		node = node.next
	}

	return nil
}

func (l *LinkedList) Len() uint32 {
	return l.len
}

func (l *LinkedList) Append(payload interface{}) {
	l.len++

	node := &LinkedNode{Payload: payload}
	if l.tail == nil {
		l.tail = node
		l.header = node
		return
	}

	l.tail.next = node
	node.pre = l.tail
	l.tail = node
}

func (l *LinkedList) Unshift(payload interface{}) {
	l.len++

	node := &LinkedNode{Payload: payload}
	if l.header == nil {
		l.header = node
		l.tail = node
		return
	}

	node.next = l.header
	l.header.pre = node
	l.header = node
}

func (l *LinkedList) Delete(payload interface{}) bool {
	node, ok := l.Search(payload)
	if !ok {
		return false
	}

	l.len--

	if node.pre != nil {
		node.pre = node.next
	}

	if node.next != nil {
		node.next.pre = node.pre
	}

	return true
}

var ErrIdxOverLinkedListRange = errors.New("index over linked list range")

func (l *LinkedList) Insert(index uint32, payload interface{}) error {
	linkLen := l.len

	if index < 0 || index > linkLen {
		return ErrIdxOverLinkedListRange
	}

	l.len++

	node := l.header
	var i uint32
	for ; i < index; i++ {
		node = node.next
	}

	newNode := &LinkedNode{Payload: payload}
	newNode.next = node.next
	newNode.pre = node
	node.next = newNode

	return nil
}

func (l *LinkedList) Search(payload interface{}) (*LinkedNode, bool) {
	var dest *LinkedNode
	_ = l.Walk(func(node *LinkedNode) (bool, error) {
		if node.Payload == payload {
			dest = node
			return false, nil
		}
		return true, nil
	})

	if dest != nil {
		return dest, true
	}

	return nil, false
}

type LinkedNode struct {
	Payload interface{}
	next    *LinkedNode
	pre     *LinkedNode
}
