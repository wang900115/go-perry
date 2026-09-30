package main

import "fmt"

type Relation int

const (
	Equal Relation = iota
	Before
	After
	Concurrent
)

func (r Relation) String() string {
	switch r {
	case Equal:
		return "Equal"
	case Before:
		return "Before"
	case After:
		return "After"
	case Concurrent:
		return "Concurrent"
	default:
		return "Unknown"
	}
}

type Path []bool

func copyPath(p Path) Path {
	result := make(Path, len(p))
	copy(result, p)
	return result
}

func equalPath(a, b Path) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

type Node struct {
	left    *Node
	right   *Node
	counter uint64
}

func cloneNode(n *Node) *Node {
	if n == nil {
		return nil
	}

	return &Node{
		left:    cloneNode(n.left),
		right:   cloneNode(n.right),
		counter: n.counter,
	}
}

type ITC struct {
	root *Node
	path Path
}

func NewITC() ITC {
	return ITC{
		root: &Node{},
	}
}

// current returns the node identified by the current path.
func (c *ITC) current() *Node {
	n := c.root

	for _, direction := range c.path {
		if direction {
			if n.right == nil {
				n.right = &Node{}
			}
			n = n.right
		} else {
			if n.left == nil {
				n.left = &Node{}
			}
			n = n.left
		}
	}

	return n
}

// Event records one local event.
func (c *ITC) Event() {
	c.current().counter++
}

// Fork splits the current identity interval.
//
// The current node becomes an internal node:
//
//             old
//             / \
//            /   \
//         child child
//
// The caller keeps the left interval.
// The returned ITC owns the right interval.
func (c *ITC) Fork() ITC {
	current := c.current()

	left := &Node{
		counter: current.counter,
	}

	right := &Node{
		counter: current.counter,
	}

	current.left = left
	current.right = right
	current.counter = 0

	// Current process keeps left.
	c.path = append(copyPath(c.path), false)

	// New process owns right.
	childPath := append(copyPath(c.path[:len(c.path)-1]), true)

	return ITC{
		root: c.root,
		path: childPath,
	}
}

// Clone creates an independent copy.
func (c ITC) Clone() ITC {
	return ITC{
		root: cloneNode(c.root),
		path: copyPath(c.path),
	}
}

// Join merges another clock into this clock.
//
// For the educational implementation we recursively take
// the maximum counter at corresponding tree positions.
func (c *ITC) Join(other ITC) {
	c.root = joinNodes(c.root, other.root)
}

func joinNodes(a, b *Node) *Node {
	if a == nil && b == nil {
		return nil
	}

	if a == nil {
		return cloneNode(b)
	}

	if b == nil {
		return cloneNode(a)
	}

	result := &Node{}

	if a.counter > b.counter {
		result.counter = a.counter
	} else {
		result.counter = b.counter
	}

	result.left = joinNodes(a.left, b.left)
	result.right = joinNodes(a.right, b.right)

	return result
}

// Compare determines the causal relationship.
//
// We compare the whole tree:
//
// Equal
// Before
// After
// Concurrent
func (c ITC) Compare(other ITC) Relation {
	aLEB, aGEB := compareNodes(c.root, other.root)

	switch {
	case aLEB && aGEB:
		return Equal
	case aLEB:
		return Before
	case aGEB:
		return After
	default:
		return Concurrent
	}
}

// compareNodes returns:
//
// aLEB = a <= b
// aGEB = a >= b
//
// If both are true -> Equal
// Only aLEB       -> Before
// Only aGEB       -> After
// Neither         -> Concurrent
func compareNodes(a, b *Node) (aLEB, aGEB bool) {
	if a == nil && b == nil {
		return true, true
	}

	if a == nil {
		if isZeroTree(b) {
			return true, true
		}
		return true, false
	}

	if b == nil {
		if isZeroTree(a) {
			return true, true
		}
		return false, true
	}

	lessEqual := a.counter <= b.counter
	greaterEqual := a.counter >= b.counter

	leftLE, leftGE := compareNodes(a.left, b.left)
	rightLE, rightGE := compareNodes(a.right, b.right)

	return lessEqual && leftLE && rightLE, greaterEqual && leftGE && rightGE
}

func isZeroTree(n *Node) bool {
	if n == nil {
		return true
	}

	if n.counter != 0 {
		return false
	}

	return isZeroTree(n.left) && isZeroTree(n.right)
}

// Debug printing.

func (c ITC) Print() {
	printNode(c.root, "", "root")
	fmt.Printf("path=%v\n", c.path)
}

func printNode(n *Node, prefix, name string) {
	if n == nil {
		return
	}

	fmt.Printf(
		"%s%s counter=%d\n",
		prefix,
		name,
		n.counter,
	)

	if n.left != nil {
		printNode(n.left, prefix+"  ", "left")
	}

	if n.right != nil {
		printNode(n.right, prefix+"  ", "right")
	}
}

func main() {
	// --------------------------------------------------
	// P1
	// --------------------------------------------------

	p1 := NewITC()

	p1.Event()
	p1.Event()
	p1.Event()

	fmt.Println("P1 after 3 events:")
	p1.Print()

	// --------------------------------------------------
	// Fork
	// --------------------------------------------------

	p2 := p1.Fork()

	fmt.Println("\nAfter fork:")
	fmt.Println("P1:")
	p1.Print()

	fmt.Println("P2:")
	p2.Print()

	// --------------------------------------------------
	// Independent events
	// --------------------------------------------------

	p1.Event()

	p2.Event()
	p2.Event()

	fmt.Println("\nAfter independent events:")

	fmt.Println("P1:")
	p1.Print()

	fmt.Println("P2:")
	p2.Print()

	// --------------------------------------------------
	// Compare
	// --------------------------------------------------

	fmt.Println("\nCompare:")
	fmt.Println("P1 vs P2:", p1.Compare(p2))
	fmt.Println("P2 vs P1:", p2.Compare(p1))

	// --------------------------------------------------
	// Join
	// --------------------------------------------------

	p1.Join(p2)

	fmt.Println("\nAfter P1 joins P2:")
	fmt.Println("P1:")
	p1.Print()

	fmt.Println("P1 vs P2:", p1.Compare(p2))

	// --------------------------------------------------
	// Clone
	// --------------------------------------------------

	p3 := p1.Clone()

	p3.Event()

	fmt.Println("\nClone:")
	fmt.Println("P1 vs P3:", p1.Compare(p3))
	fmt.Println("P3 vs P1:", p3.Compare(p1))
}
