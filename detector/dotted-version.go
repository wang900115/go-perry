package main

import (
	"fmt"
	"sort"
)

type Actor string

// VersionVector records the causal history that has been observed.
type VersionVector map[Actor]uint64

// Dot identifies exactly one event:
// actor A's N-th event.
type Dot struct {
	Actor   Actor
	Counter uint64
}

// DVV = causal context + one specific event.
type DVV struct {
	Context VersionVector
	Dot     Dot
}

// Next creates a new dot for actor.
func (vv VersionVector) Next(actor Actor) Dot {
	return Dot{
		Actor:   actor,
		Counter: vv[actor] + 1,
	}
}

// Advance adds a dot into the version vector.
func (vv VersionVector) Advance(dot Dot) {
	if dot.Counter > vv[dot.Actor] {
		vv[dot.Actor] = dot.Counter
	}
}

// NewDVV creates a DVV from an existing causal context.
func NewDVV(context VersionVector, actor Actor) DVV {
	dot := context.Next(actor)

	return DVV{
		Context: cloneVV(context),
		Dot:     dot,
	}
}

// Merge combines two DVVs into a causal context.
func Merge(a, b DVV) VersionVector {
	result := cloneVV(a.Context)

	result.Advance(a.Dot)
	result.Advance(b.Dot)

	for actor, counter := range b.Context {
		if counter > result[actor] {
			result[actor] = counter
		}
	}

	return result
}

// Contains checks whether the version vector has already seen a dot.
func Contains(vv VersionVector, dot Dot) bool {
	return vv[dot.Actor] >= dot.Counter
}

// HappensBefore determines whether a DVV event happened before another.
//
// a < b means:
//   - a's dot is contained in b's causal history
//   - and b has additional causal information.
func HappensBefore(a, b DVV) bool {
	// a's event must be known by b.
	if !Contains(b.Context, a.Dot) {
		return false
	}

	// b must contain something that a does not know.
	bv := Merge(b, DVV{})

	av := Merge(a, DVV{})

	return !equalVV(av, bv)
}

func cloneVV(vv VersionVector) VersionVector {
	result := make(VersionVector, len(vv))

	for actor, counter := range vv {
		result[actor] = counter
	}

	return result
}

func equalVV(a, b VersionVector) bool {
	if len(a) != len(b) {
		return false
	}

	for actor, counter := range a {
		if b[actor] != counter {
			return false
		}
	}

	return true
}

func printVV(vv VersionVector) {
	actors := make([]string, 0, len(vv))

	for actor := range vv {
		actors = append(actors, string(actor))
	}

	sort.Strings(actors)

	fmt.Print("{")

	for i, actor := range actors {
		if i > 0 {
			fmt.Print(", ")
		}

		fmt.Printf("%s:%d", actor, vv[Actor(actor)])
	}

	fmt.Println("}")
}

func main() {
	// Initial causal history.
	vv := VersionVector{}

	// Replica A creates its first event.
	a1 := NewDVV(vv, "A")

	fmt.Println("A1")
	fmt.Print("Context = ")
	printVV(a1.Context)
	fmt.Printf("Dot = (%s,%d)\n\n",
		a1.Dot.Actor,
		a1.Dot.Counter,
	)

	// A1 becomes part of A's causal history.
	vv.Advance(a1.Dot)

	// Replica A creates another event.
	a2 := NewDVV(vv, "A")

	fmt.Println("A2")
	fmt.Print("Context = ")
	printVV(a2.Context)
	fmt.Printf("Dot = (%s,%d)\n\n",
		a2.Dot.Actor,
		a2.Dot.Counter,
	)

	// Replica B receives A2's history and creates B1.
	bContext := Merge(a2, DVV{})

	b1 := NewDVV(bContext, "B")

	fmt.Println("B1")
	fmt.Print("Context = ")
	printVV(b1.Context)
	fmt.Printf("Dot = (%s,%d)\n\n",
		b1.Dot.Actor,
		b1.Dot.Counter,
	)

	// Replica A independently creates A3 from A2.
	a3 := NewDVV(a2.Context, "A")

	fmt.Println("A3")
	fmt.Print("Context = ")
	printVV(a3.Context)
	fmt.Printf("Dot = (%s,%d)\n\n",
		a3.Dot.Actor,
		a3.Dot.Counter,
	)

	fmt.Println("Causal relationships:")

	fmt.Println("A2 -> B1:", HappensBefore(a2, b1))
	fmt.Println("A2 -> A3:", HappensBefore(a2, a3))
	fmt.Println("B1 -> A3:", HappensBefore(b1, a3))
}
