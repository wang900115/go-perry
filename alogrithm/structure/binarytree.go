package binarytree

type Node struct {
	val   int
	left  *Node
	right *Node
}

type BinaryTree struct {
	root *Node
}

func (b *BinaryTree) Insert(val int) bool {
	if b.root == nil {
		b.root = &Node{val: val}
		return true
	}
	curr := b.root
	for curr.val != val {
		if val > curr.val {
			if curr.right == nil {
				curr.right = &Node{val: val}
				return true
			}
			curr = curr.right
		} else {
			if curr.left == nil {
				curr.left = &Node{val: val}
				return true
			}
			curr = curr.left
		}
	}
	return false
}

func (b *BinaryTree) Search(val int) bool {
	curr := b.root
	for curr != nil {
		if val == curr.val {
			return true
		} else if val < curr.val {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return false
}

// * Important
func (b *BinaryTree) Delete(val int) bool {
	var parent *Node
	curr := b.root

	for curr != nil && curr.val != val {
		if curr.val > val {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	if curr == nil {
		return false
	}

	var child *Node
	if curr.left == nil {
		child = curr.right
	} else if curr.right == nil {
		child = curr.left
	} else {
		successorParent := curr
		successor := curr.right
		for successor.left != nil {
			successorParent = successor
			successor = successor.left
		}

		curr.val = successor.val
		if successorParent.left == successor {
			successorParent.right = successor.right
		} else {
			successorParent.left = successor.right
		}
		return true
	}

	if parent == nil {
		b.root = child
	} else if parent.left == curr {
		parent.left = child
	} else {
		parent.right = child
	}
	return true
}

func (b *BinaryTree) Min() (int, bool) {
	if b.root == nil {
		return 0, false
	}
	curr := b.root
	for curr.left != nil {
		curr = curr.left
	}
	return curr.val, true
}

func (b *BinaryTree) Max() (int, bool) {
	if b.root == nil {
		return 0, false
	}
	curr := b.root
	for curr.right != nil {
		curr = curr.right
	}
	return curr.val, true
}

func (b *BinaryTree) InOrder() {
}

func (b *BinaryTree) PreOrder() {
}

func (b *BinaryTree) PostOrder() {
}
