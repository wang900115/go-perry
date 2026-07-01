package main

// bounds check elimination: the compiler can eliminate bounds checks for array/slice accesses when it can prove that the index is within bounds
// the compiler can eliminate bounds checks for array/slice accesses when it can prove that the index is within bounds

func sum(a []int) int {
	sum := 0

	for i := 0; i < len(a); i++ {
		sum += a[i]
	}

	return sum
}

// need check every time, because the compiler cannot prove that the index is within bounds
func sum2(a []int, index []int) int {
	sum := 0

	for _, i := range index {
		sum += a[i]
	}

	return sum
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	println(sum(s))
	println(sum2(s, []int{0, 1, 2, 3, 4}))
}
