package main

import (
	"fmt"
	"sort"
)

type NodeID string

type VersionVector map[NodeID]uint64

// Increment creates a new local event.
func (vv VersionVector) Increment(node NodeID) {
	vv[node]++
}

// Update merges another version vector.
// For every node, keep the maximum counter.
func (vv VersionVector) Update(other VersionVector) {
	for node, counter := range other {
		if counter > vv[node] {
			vv[node] = counter
		}
	}
}

// Clone returns a copy of the version vector.
func (vv VersionVector) Clone() VersionVector {
	result := make(VersionVector, len(vv))

	for node, counter := range vv {
		result[node] = counter
	}

	return result
}

// HappensBefore returns true if vv < other.
func (vv VersionVector) HappensBefore(other VersionVector) bool {
	less := false

	nodes := make(map[NodeID]struct{})

	for node := range vv {
		nodes[node] = struct{}{}
	}

	for node := range other {
		nodes[node] = struct{}{}
	}

	for node := range nodes {
		a := vv[node]
		b := other[node]

		if a > b {
			return false
		}

		if a < b {
			less = true
		}
	}

	return less
}

// Equal checks whether two version vectors are identical.
func (vv VersionVector) Equal(other VersionVector) bool {
	return !vv.HappensBefore(other) &&
		!other.HappensBefore(vv)
}

// Concurrent returns true when neither vector happened-before the other.
func (vv VersionVector) Concurrent(other VersionVector) bool {
	return !vv.HappensBefore(other) &&
		!other.HappensBefore(vv) &&
		!vv.Equal(other)
}

func (vv VersionVector) String() string {
	nodes := make([]string, 0, len(vv))

	for node := range vv {
		nodes = append(nodes, string(node))
	}

	sort.Strings(nodes)

	result := "{"

	for i, node := range nodes {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("%s:%d", node, vv[NodeID(node)])
	}

	result += "}"

	return result
}

func main() {
	// ----------------------------------------
	// Node A creates an event
	// ----------------------------------------

	a := VersionVector{}

	a.Increment("A")

	fmt.Println("A:", a)
	// A: {A:1}

	// ----------------------------------------
	// Node A creates another event
	// ----------------------------------------

	a.Increment("A")

	fmt.Println("A:", a)
	// A: {A:2}

	// ----------------------------------------
	// Node B receives A's state
	// ----------------------------------------

	b := a.Clone()

	b.Increment("B")

	fmt.Println("B:", b)
	// B: {A:2, B:1}

	// ----------------------------------------
	// Node A creates another event independently
	// ----------------------------------------

	a.Increment("A")

	fmt.Println("A:", a)
	// A: {A:3}

	// ----------------------------------------
	// Compare A and B
	// ----------------------------------------

	fmt.Println("A happens-before B:", a.HappensBefore(b))
	fmt.Println("B happens-before A:", b.HappensBefore(a))
	fmt.Println("A concurrent B:", a.Concurrent(b))
}
