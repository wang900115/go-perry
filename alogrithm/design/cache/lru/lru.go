package lru

type LRUCache struct {
	cache    map[int]*Node
	capacity int
	head     *Node
	tail     *Node
}

type Node struct {
	key, val   int
	prev, next *Node
}

func Constructor(capacity int) LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head
	return LRUCache{
		cache:    make(map[int]*Node),
		capacity: capacity,
		head:     head,
		tail:     tail,
	}
}

func (this *LRUCache) Get(key int) int {
	if node, ok := this.cache[key]; ok {
		this.moveToHead(node)
		return node.val
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if node, ok := this.cache[key]; ok {
		node.val = value
		this.moveToHead(node)
		return
	}

	if len(this.cache) == this.capacity {
		lru := this.removeTail()
		delete(this.cache, lru.key)
	}

	newNode := &Node{key: key, val: value}
	this.cache[key] = newNode
	this.moveToHead(newNode)
}

func (this *LRUCache) moveToHead(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.next = this.head.next
	node.prev = this.head
	this.head.next.prev = node
	this.head.next = node
}

func (this *LRUCache) removeTail() *Node {
	tail := this.tail.prev
	tail.prev.next = this.tail
	this.tail.prev = tail.prev
	return tail
}

func (this *LRUCache) PrintCache() {
	curr := this.head.next
	for curr != this.tail {
		println("Key:", curr.key, "Value:", curr.val)
		curr = curr.next
	}
}

func main() {
	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	println(cache.Get(1))
	cache.Put(3, 3)
	cache.PrintCache()
	println(cache.Get(2))
	cache.Put(4, 4)
	cache.PrintCache()
	println(cache.Get(1))
	println(cache.Get(3))
	println(cache.Get(4))
}
