package char

type SteamCheck struct {
	root   *Node
	stream []byte
	maxLen int
}

type Node struct {
	childrens [26]*Node
	isWord    bool
}

func Constructor(words []string) SteamCheck {
	root := &Node{}
	maxLen := 0
	for _, word := range words {
		if len(word) > maxLen {
			maxLen = len(word)
		}
		inversedInsert(root, word)
	}
	return SteamCheck{
		root:   root,
		maxLen: maxLen,
	}
}

func inversedInsert(root *Node, word string) {
	node := root
	for i := len(word) - 1; i >= 0; i-- {
		index := word[i] - 'a'
		if node.childrens[index] == nil {
			node.childrens[index] = &Node{}
		}
		node = node.childrens[index]
	}
	node.isWord = true
}

func (this *SteamCheck) Query(letter byte) bool {
	this.stream = append(this.stream, letter)
	if len(this.stream) > this.maxLen {
		this.stream = this.stream[1:]
	}
	node := this.root
	for i := len(this.stream) - 1; i >= 0; i-- {
		index := this.stream[i] - 'a'
		if node.childrens[index] == nil {
			return false
		}
		node = node.childrens[index]
		if node.isWord {
			return true
		}
	}
	return false
}
