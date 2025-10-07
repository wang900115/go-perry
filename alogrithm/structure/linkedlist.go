package linkedlist

type Node struct {
	val  int
	next *Node
}

type LinkedList struct {
	head *Node
}

func (l *LinkedList) Append(val int) {
	newNode := &Node{val: val}
	if l.head == nil {
		l.head = newNode
		return
	}

	current := l.head
	for current.next != nil {
		current = current.next
	}
	current.next = newNode
}

func (l *LinkedList) Prepend(val int) {
	newNode := &Node{val: val}
	newNode.next = l.head.next
	l.head = newNode
}

func (l *LinkedList) Delete(val int) {
	if l.head == nil {
		return
	}
	if l.head.val == val {
		l.head = l.head.next
		return
	}
	current := l.head
	for current.next != nil && current.next.val != val {
		current = current.next
	}
	if current.next != nil {
		current.next = current.next.next
	}
}

func (l *LinkedList) Find(val int) bool {
	current := l.head
	for current != nil {
		if current.val == val {
			return true
		}
		current = current.next
	}
	return false
}
