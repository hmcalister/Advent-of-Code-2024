package main

type LinkedList struct {
	head *LinkedListNode
	tail *LinkedListNode
}

func NewDiskMapLinkedList() *LinkedList {
	return &LinkedList{
		head: nil,
		tail: nil,
	}
}

func (list *LinkedList) Add(file *FileInformation) {
	newNode := &LinkedListNode{
		fileInfo: file,
	}
	if list.head == nil {
		list.head = newNode
	} else {
		list.tail.next = newNode
		newNode.prev = list.tail
	}
	list.tail = newNode
}

func (list *LinkedList) SpliceOut(node *LinkedListNode) {
	if node == list.head && node == list.tail {
		list.head = nil
		list.tail = nil
		return
	}

	if node == list.head {
		list.head = node.next
		list.head.prev = nil
		return
	}

	if node == list.tail {
		list.tail = node.prev
		list.tail.next = nil
		return
	}

	previousNode := node.prev
	nextNode := node.next
	previousNode.next = nextNode
	nextNode.prev = previousNode
}

func (list *LinkedList) SpliceIn(previousNode, newNode *LinkedListNode) {
	if previousNode == nil {
		newNode.next = list.head
		newNode.prev = nil
		list.head.prev = newNode
		list.head = newNode
		return
	}

	if previousNode.next == nil {
		newNode.next = nil
		newNode.prev = previousNode
		previousNode.next = newNode
		return
	}

	newNode.next = previousNode.next
	newNode.prev = previousNode
	previousNode.next = newNode
	newNode.next.prev = newNode

}

type LinkedListNode struct {
	next *LinkedListNode
	prev *LinkedListNode

	fileInfo *FileInformation
}
