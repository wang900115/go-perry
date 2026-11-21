// Package main demonstrates Go generics: type parameters, constraints,
// type inference, generic types, and the pitfalls around them.
package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ---------------------------------------------------------------------------
// Part 1: generic functions and type inference
// ---------------------------------------------------------------------------

// Map applies f to every element. Two type parameters: In and Out.
func Map[In, Out any](s []In, f func(In) Out) []Out {
	out := make([]Out, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

// Filter keeps the elements for which keep returns true.
// S ~[]E keeps the caller's named slice type in the result.
func Filter[S ~[]E, E any](s S, keep func(E) bool) S {
	var out S
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Part 2: constraints are interfaces (type sets)
// ---------------------------------------------------------------------------

// Number is a constraint: an interface with a union of types.
// The ~ means "any type whose underlying type is ...".
type Number interface {
	~int | ~int64 | ~float64
}

// Sum works for every type in Number's type set.
func Sum[T Number](xs []T) T {
	var total T // zero value of T
	for _, x := range xs {
		total += x // allowed: every type in the set supports +
	}
	return total
}

// Celsius is a named type; its underlying type is float64.
type Celsius float64

// StrictNumber has no ~, so Celsius is NOT in its type set.
type StrictNumber interface {
	int | float64
}

func SumStrict[T StrictNumber](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// MaxOf uses cmp.Ordered from the standard library (supports < <= > >=).
func MaxOf[T cmp.Ordered](first T, rest ...T) T {
	m := first
	for _, v := range rest {
		if v > m {
			m = v
		}
	}
	return m
}

// Index uses comparable (supports == and !=).
func Index[T comparable](s []T, target T) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

// Stringish mixes a method requirement with a type set.
type Stringish interface {
	~string
	Len() int
}

type Name string

func (n Name) Len() int { return len(n) }

func Longest[T Stringish](xs []T) T {
	var best T
	for _, x := range xs {
		if x.Len() > best.Len() {
			best = x
		}
	}
	return best
}

// ---------------------------------------------------------------------------
// Part 3: generic types
// ---------------------------------------------------------------------------

// Stack is a generic type. Its methods use the type's parameter T,
// but a method cannot declare new type parameters of its own.
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop returns the zero value of T and false when the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items[len(s.items)-1] = zero // do not keep a reference alive
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// Not allowed: methods cannot have type parameters.
// func (s *Stack[T]) Map[U any](f func(T) U) *Stack[U] { ... }
// -> syntax error: method must have no type parameters

// Pair shows a generic type with two parameters.
type Pair[K comparable, V any] struct {
	Key K
	Val V
}

func (p Pair[K, V]) String() string { return fmt.Sprintf("%v=%v", p.Key, p.Val) }

// ---------------------------------------------------------------------------
// Part 4: limits of type inference
// ---------------------------------------------------------------------------

// Zero has T only in the result, so T cannot be inferred from arguments.
func Zero[T any]() T {
	var z T
	return z
}

// Describe shows how to branch on the dynamic type: convert to any first.
func Describe[T any](v T) string {
	switch x := any(v).(type) {
	case int:
		return fmt.Sprintf("int %d", x)
	case string:
		return fmt.Sprintf("string %q", x)
	default:
		return fmt.Sprintf("other %T", x)
	}
}

// ---------------------------------------------------------------------------
// Part 5: comparable vs runtime panics
// ---------------------------------------------------------------------------

// Since Go 1.20, interface types like any satisfy comparable.
// The comparison then happens at run time and can panic.
func safeIndex[T comparable](s []T, target T) (idx int, err error) {
	defer func() {
		if r := recover(); r != nil {
			idx, err = -1, fmt.Errorf("recovered: %v", r)
		}
	}()
	return Index(s, target), nil
}

// ---------------------------------------------------------------------------
// Part 6: self-referential constraints (Go 1.26)
// ---------------------------------------------------------------------------

// Lesser is a generic interface that refers to itself in its own
// type parameter list. This was a compile error before Go 1.26.
type Lesser[T Lesser[T]] interface {
	Less(T) bool
}

type Version struct{ Major, Minor int }

func (a Version) Less(b Version) bool {
	if a.Major != b.Major {
		return a.Major < b.Major
	}
	return a.Minor < b.Minor
}

// MinOf works for any T that can compare itself with another T.
func MinOf[T Lesser[T]](xs []T) (T, error) {
	var zero T
	if len(xs) == 0 {
		return zero, errors.New("empty slice")
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x.Less(m) {
			m = x
		}
	}
	return m, nil
}

func main() {
	fmt.Println("== Part 1: generic functions and inference ==")
	words := []string{"go", "generics", "are", "fun"}
	lengths := Map(words, func(s string) int { return len(s) }) // In=string, Out=int inferred
	fmt.Println("lengths     =", lengths)
	upper := Map[string, string](words, strings.ToUpper) // explicit instantiation
	fmt.Println("upper       =", upper)

	type Words []string
	long := Filter(Words(words), func(s string) bool { return len(s) > 2 })
	fmt.Printf("long        = %v (type %T)\n", long, long)

	fmt.Println("\n== Part 2: constraints ==")
	fmt.Println("Sum ints    =", Sum([]int{1, 2, 3}))
	fmt.Println("Sum floats  =", Sum([]float64{1.5, 2.25}))
	temps := []Celsius{20.5, 21.5}
	fmt.Printf("Sum Celsius = %v (type %T)\n", Sum(temps), Sum(temps))
	// SumStrict(temps) // Celsius does not satisfy StrictNumber (possibly missing ~ for float64 in StrictNumber)
	fmt.Println("SumStrict   =", SumStrict([]int{4, 5}))
	fmt.Println("MaxOf       =", MaxOf(3, 9, 4), MaxOf("pear", "apple", "zoo"))
	fmt.Println("Index       =", Index([]string{"a", "b", "c"}, "c"))
	fmt.Println("Longest     =", Longest([]Name{"ann", "perry", "bob"}))
	// Longest([]string{"a"}) // string does not satisfy Stringish (missing method Len)

	fmt.Println("\n== Part 3: generic types ==")
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	v, ok := s.Pop()
	fmt.Println("pop         =", v, ok)
	v, ok = s.Pop()
	fmt.Println("pop         =", v, ok)
	v, ok = s.Pop()
	fmt.Println("pop empty   =", v, ok)
	ps := []Pair[string, int]{{"b", 2}, {"a", 1}, {"c", 3}}
	slices.SortFunc(ps, func(x, y Pair[string, int]) int { return cmp.Compare(x.Key, y.Key) })
	fmt.Println("sorted pairs=", ps)

	fmt.Println("\n== Part 4: inference limits ==")
	// z := Zero() // cannot infer T
	fmt.Printf("Zero[string]= %q, Zero[*int] = %v\n", Zero[string](), Zero[*int]())
	// f := Sum // cannot use generic function Sum without instantiation
	f := Sum[float64]
	fmt.Printf("f is %T, f(...) = %v\n", f, f([]float64{0.5, 0.5}))
	fmt.Println(Describe(42), "|", Describe("hi"), "|", Describe(3.14))

	fmt.Println("\n== Part 5: comparable and runtime panics ==")
	idx, err := safeIndex([]any{1, "x", 2.5}, any("x"))
	fmt.Println("[]any ok    =", idx, err)
	idx, err = safeIndex([]any{[]int{1}}, any([]int{1}))
	fmt.Println("[]any slice =", idx, err)

	fmt.Println("\n== Part 6: self-referential constraint (Go 1.26) ==")
	m, err := MinOf([]Version{{1, 22}, {1, 9}, {1, 26}})
	fmt.Println("MinOf       =", m, err)
	_, err = MinOf([]Version{})
	fmt.Println("MinOf empty =", err)
}
