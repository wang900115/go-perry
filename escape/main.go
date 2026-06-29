package main

import "fmt"

// can inline NewUser function to avoid heap allocation
// inlining call to NewUser function will avoid heap allocation and improve performance
// escape analysis: the NewUser function does not escape to the heap, so it can be inlined
// if escape leads to heap allocation, it can cause performance issues(lower cache efficiency) and increase memory usage
// three way escape normals: 1. return pointer 2. interface 3. closure
type User struct {
	Name string
}

func NewUser(name string) *User {
	u := &User{Name: name}
	return u
}

func main() {
	u := NewUser("Perry")
	fmt.Println(u.Name)
}

// HPC core:
// 1. using value instead of pointer to avoid heap allocation, if need for cross scope, use pointer to avoid copying large struct.
// 2. avoid return pointer if just data.
// 3. using sync.Pool to reuse objects.
// 4. preallocate memory for large slices to avoid multiple allocations and reduce GC overhead.
// 5. avoid using closure capture with goroutines.
// 6. avoid fmt, using structed logging or buffer reuse.
// 7. using benchmark test to find hot spots and optimize them.
