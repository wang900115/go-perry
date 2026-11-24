// Package main demonstrates how Go maps behave: nil maps, reference-like
// sharing, comma-ok lookups, random iteration order, deterministic iteration
// with the maps/slices packages, non-addressable elements, mutation during
// range, NaN keys, delete vs clear, and the concurrent-write fatal error.
package main

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"strings"
)

type Point struct{ X, Y int }

func main() {
	section("1. nil map: reads are fine, writes panic")
	var nilMap map[string]int
	fmt.Println("len:", len(nilMap), "lookup:", nilMap["missing"], "nil?", nilMap == nil)
	for k := range nilMap { // ranging over nil is a no-op
		fmt.Println("never printed", k)
	}
	fmt.Println("write err:", tryWrite(nilMap))

	section("2. a map value is a pointer to the runtime map")
	a := map[string]int{"x": 1}
	b := a // copies the pointer, not the entries
	b["y"] = 2
	fmt.Println("a:", a, "b:", b)
	addEntry(a)
	fmt.Println("after addEntry(a):", a)
	c := maps.Clone(a) // shallow copy: a new map with the same entries
	c["x"] = 100
	fmt.Println("a after editing clone:", a, "clone:", c)

	section("3. comma-ok: zero value vs missing key")
	stock := map[string]int{"apple": 0}
	v1, ok1 := stock["apple"]
	v2, ok2 := stock["banana"]
	fmt.Printf("apple=%d ok=%v | banana=%d ok=%v\n", v1, ok1, v2, ok2)

	section("4. iteration order is unspecified (and deliberately randomized)")
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8}
	seen := map[string]bool{}
	for range 50 {
		var sb strings.Builder
		for k := range m {
			sb.WriteString(k)
		}
		seen[sb.String()] = true
	}
	fmt.Println("distinct orders in 50 runs > 1:", len(seen) > 1)
	keys := slices.Sorted(maps.Keys(m)) // deterministic order
	fmt.Println("sorted keys:", keys)
	fmt.Println("fmt prints maps sorted:", map[int]string{3: "c", 1: "a", 2: "b"})

	section("5. map elements are not addressable")
	pts := map[string]Point{"p": {1, 2}}
	// pts["p"].X = 10 // compile error: cannot assign to struct field pts["p"].X in map
	p := pts["p"] // copy out, modify, store back
	p.X = 10
	pts["p"] = p
	fmt.Println("value map after store-back:", pts)
	ptrs := map[string]*Point{"p": {1, 2}}
	ptrs["p"].X = 10 // fine: the element is a pointer, we modify what it points to
	fmt.Println("pointer map:", *ptrs["p"])
	counts := map[string]int{}
	for _, w := range strings.Fields("go is fun go is fast go") {
		counts[w]++ // m[k]++ works: it is m[k] = m[k] + 1
	}
	fmt.Println("word counts:", counts)

	section("6. deleting during range is safe; adding is unpredictable")
	nums := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true}
	for k := range nums {
		if k%2 == 0 {
			delete(nums, k)
		}
	}
	fmt.Println("odd keys left:", slices.Sorted(maps.Keys(nums)))
	grow := map[int]int{0: 0}
	visited := 0
	for k := range grow {
		visited++
		if k < 100 {
			grow[k+1] = k + 1 // may or may not be visited in this loop
		}
	}
	fmt.Println("visited between 1 and 101 (exact count varies):", visited >= 1 && visited <= 101)

	section("7. NaN keys: you can insert, but never find or delete them")
	f := map[float64]string{}
	nan := math.NaN()
	f[nan] = "first"
	f[nan] = "second" // NaN != NaN, so this is a NEW entry
	_, found := f[nan]
	delete(f, nan) // no-op for NaN
	fmt.Println("len:", len(f), "found:", found)
	clear(f) // Go 1.21: clear removes everything, including NaN keys
	fmt.Println("len after clear:", len(f))

	section("8. set idiom and interface keys")
	set := map[string]struct{}{}
	for _, s := range []string{"go", "rust", "go"} {
		set[s] = struct{}{}
	}
	_, hasGo := set["go"]
	fmt.Println("set size:", len(set), "has go:", hasGo)
	fmt.Println("unhashable interface key err:", tryInterfaceKey())

	section("9. concurrent writes")
	fmt.Println("concurrent map writes are a fatal error, not a panic: recover() cannot catch it.")
	fmt.Println("use sync.Mutex/RWMutex or sync.Map (see the sync_map topic).")
}

func section(title string) { fmt.Println("\n==", title) }

func addEntry(m map[string]int) { m["fromFunc"] = 3 }

func tryWrite(m map[string]int) (err any) {
	defer func() { err = recover() }()
	m["boom"] = 1
	return nil
}

func tryInterfaceKey() (err any) {
	defer func() { err = recover() }()
	m := map[any]int{}
	m[[]int{1}] = 1 // slices are not comparable: run-time panic
	return nil
}
