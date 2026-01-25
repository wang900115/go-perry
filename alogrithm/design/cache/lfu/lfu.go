package lfu

type LFUCache struct {
	nodes    map[int]*Node
	freqMap  map[int]*DLinkedList
	minFreq  int
	capacity int
	size     int
}

type Node struct {
	key, val, freq int
	prev, next     *Node
}

type DLinkedList struct {
	head, tail *Node
	size       int
}

func NewDLinkedList() *DLinkedList {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head
	return &DLinkedList{
		head: head,
		tail: tail,
		size: 0,
	}
}

func (d *DLinkedList) addToHead(node *Node) {
	node.next = d.head.next
	node.prev = d.head
	d.head.next.prev = node
	d.head.next = node
	d.size++
}

func (d *DLinkedList) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
	d.size--
}

func (d *DLinkedList) removeTail() *Node {
	if d.size == 0 {
		return nil
	}
	tailNode := d.tail.prev
	d.removeNode(tailNode)
	return tailNode
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		nodes:    make(map[int]*Node),
		freqMap:  make(map[int]*DLinkedList),
		capacity: capacity,
		minFreq:  0,
		size:     0,
	}
}

func (this *LFUCache) Get(key int) int {
	if node, exists := this.nodes[key]; exists {
		this.updateNodeFreq(node)
		return node.val
	}
	return -1
}

func (this *LFUCache) Put(key int, value int) {
	if this.capacity == 0 {
		return
	}
	if node, exists := this.nodes[key]; exists {
		node.val = value
		this.updateNodeFreq(node)
		return
	}
	if this.size == this.capacity {
		minFreqList := this.freqMap[this.minFreq]
		removedNode := minFreqList.removeTail()
		delete(this.nodes, removedNode.key)
		this.size--
	}
	newNode := &Node{key: key, val: value, freq: 1}
	this.nodes[key] = newNode
	if _, exists := this.freqMap[1]; !exists {
		this.freqMap[1] = NewDLinkedList()
	}
	this.freqMap[1].addToHead(newNode)
	this.minFreq = 1
	this.size++
}

func (this *LFUCache) updateNodeFreq(node *Node) {
	oldFreq := node.freq
	this.freqMap[oldFreq].removeNode(node)
	if oldFreq == this.minFreq && this.freqMap[oldFreq].size == 0 {
		this.minFreq++
	}
	node.freq++
	if _, exists := this.freqMap[node.freq]; !exists {
		this.freqMap[node.freq] = NewDLinkedList()
	}
	this.freqMap[node.freq].addToHead(node)
}
