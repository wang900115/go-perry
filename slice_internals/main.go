// Package main demonstrates how Go slices work: the slice header (pointer,
// len, cap), shared backing arrays, append aliasing, growth, passing slices
// to functions, slices.Delete zeroing, and memory retention by sub-slices.
package main

import (
	"fmt"
	"slices"
)

func header(name string, s []int) {
	fmt.Printf("%-6s len=%d cap=%d %v\n", name, len(s), cap(s), s)
}

// part1: sub-slices share the same backing array.
func part1() {
	fmt.Println("== Part 1: slice header and shared backing array ==")
	arr := [6]int{0, 1, 2, 3, 4, 5}
	a := arr[1:4] // len 3, cap 5 (from index 1 to the end of arr)
	b := a[1:2]   // len 1, cap 4
	header("a", a)
	header("b", b)
	b[0] = 99 // writes arr[2]
	fmt.Println("after b[0] = 99:")
	header("a", a)
	fmt.Println("arr  ", arr)
	// Re-slicing can go past len, up to cap.
	header("b[:4]", b[:4])
	fmt.Println()
}

// part2: two appends to the same base slice overwrite each other.
func part2() {
	fmt.Println("== Part 2: append aliasing ==")
	base := make([]int, 3, 10)
	base[0], base[1], base[2] = 1, 2, 3
	x := append(base, 100) // fits in cap: writes base's backing array at index 3
	y := append(base, 200) // same index 3 again: overwrites x[3]
	header("x", x)
	header("y", y)
	fmt.Println("x[3] changed by appending to y:", x[3] == 200)

	// Fix: a full slice expression limits cap, so append must copy.
	limited := base[:3:3]
	x2 := append(limited, 100)
	y2 := append(limited, 200)
	header("x2", x2)
	header("y2", y2)
	// slices.Clip(s) is the same as s[:len(s):len(s)].
	fmt.Println("cap(slices.Clip(base)) =", cap(slices.Clip(base)))
	fmt.Println()
}

// part3: append grows the backing array when cap is not enough.
func part3() {
	fmt.Println("== Part 3: growth (cap changes while appending 1 by 1) ==")
	var s []int
	prevCap := cap(s)
	var caps []int
	for i := 0; i < 2000; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			caps = append(caps, cap(s))
			prevCap = cap(s)
		}
	}
	fmt.Println("caps:", caps)

	// Growth moves data to a new array; the old slice keeps the old one.
	old := []int{1, 2, 3} // len 3, cap 3
	grown := append(old, 4)
	grown[0] = 42
	header("old", old)
	header("grown", grown)

	// Preallocate when you know the size: one allocation, no copies.
	pre := make([]int, 0, 2000)
	for i := 0; i < 2000; i++ {
		pre = append(pre, i)
	}
	fmt.Println("preallocated: len", len(pre), "cap", cap(pre))
	fmt.Println()
}

func setFirst(s []int)        { s[0] = -1 }
func appendOne(s []int)       { s = append(s, 7); _ = s }
func appendRet(s []int) []int { return append(s, 7) }

// part4: a slice is passed by value (the header is copied).
func part4() {
	fmt.Println("== Part 4: passing slices to functions ==")
	s := make([]int, 2, 4)
	setFirst(s)
	header("set", s)
	appendOne(s)
	header("app", s)
	fmt.Println("hidden element still in the array:", s[:3])
	s = appendRet(s)
	header("ret", s)
	fmt.Println()
}

// part5: since Go 1.22, slices.Delete zeroes the removed tail.
func part5() {
	fmt.Println("== Part 5: slices.Delete / Compact zero the tail (Go 1.22+) ==")
	s := []int{1, 2, 3, 4, 5}
	d := slices.Delete(s, 1, 3) // remove s[1:3]
	header("d", d)
	header("s", s) // old header still has len 5: sees the zeroed tail

	type user struct{ name string }
	users := []*user{{"ann"}, {"bob"}, {"cid"}}
	users2 := slices.Delete(users, 0, 1)
	fmt.Println("len(users2) =", len(users2), "users[2] == nil:", users[2] == nil, "(no leaked pointer)")
	fmt.Println()
}

// part6: a small sub-slice keeps the whole large array alive.
func part6() {
	fmt.Println("== Part 6: memory retention by sub-slices ==")
	big := make([]byte, 1<<20) // 1 MiB
	head := big[:16]
	fmt.Printf("head: len=%d cap=%d (pins the 1 MiB array)\n", len(head), cap(head))
	cloned := slices.Clone(head)
	fmt.Printf("clone: len=%d cap=%d (new small array, big can be collected)\n", len(cloned), cap(cloned))

	var nilS []int
	empty := []int{}
	fmt.Println("nil slice:", nilS == nil, len(nilS), "| empty slice:", empty == nil, len(empty))
	nilS = append(nilS, 1) // append works on a nil slice
	fmt.Println("append to nil:", nilS)
}

func main() {
	part1()
	part2()
	part3()
	part4()
	part5()
	part6()
}
