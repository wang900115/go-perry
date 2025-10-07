package main

import (
	"fmt"
	"sort"
)

func main() {
	strs := []string{"quick", "brown", "fox", "jumps"}
	// sort.Strings(strs)
	sort.Sort(sort.StringSlice(strs))
	fmt.Println("Sorted string:", strs)
	sort.Sort(sort.Reverse(sort.StringSlice(strs)))
	fmt.Println("Reversed:", strs)
	sort.Slice(strs, func(i, j int) bool {
		return len(strs[i]) < len(strs[j])
	})
	fmt.Println("Sorted String by length:", strs)
	sort.SliceStable(strs, func(i, j int) bool {
		return len(strs[i]) < len(strs[j])
	})
	fmt.Println("[Stable] Sorted String by length:", strs)
	sort.SliceStable(strs, func(i, j int) bool {
		return len(strs[j]) < len(strs[i])
	})
	fmt.Println("[Stable] Reversd Sorted String by length:", strs)

	ints := []int{56, 19, 78, 67, 14, 25}
	// sort.Ints(ints)
	sort.Sort(sort.IntSlice(ints))
	fmt.Println("Sorted int:", ints)
	sort.Sort(sort.Reverse(sort.IntSlice(ints)))
	fmt.Println("Reversed:", ints)

	floats := []float64{176.8, 19.5, 20.8, 57.4}
	// sort.Float64s(floats)
	sort.Sort(sort.Float64Slice(floats))
	fmt.Println("Sorted float:", floats)
	sort.Sort(sort.Reverse(sort.Float64Slice(floats)))
	fmt.Println("Reversed:", floats)

}
