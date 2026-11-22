// Closures and loop variables in Go.
//
// Run: go run .
// Escape analysis: go build -gcflags=-m . 2>&1 | grep -E "moved to heap|func literal"
package main

import (
	"fmt"
	"sort"
	"sync"
)

// counter returns a closure that owns its own "n".
// n outlives the call to counter, so the compiler moves it to the heap.
func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// 1. A closure captures the VARIABLE, not a snapshot of its value.
func captureByReference() {
	fmt.Println("== 1. capture by reference ==")
	x := 1
	show := func() { fmt.Println("  x =", x) }
	x = 2
	show() // prints 2, not 1

	c1, c2 := counter(), counter()
	fmt.Println("  c1:", c1(), c1(), c1())
	fmt.Println("  c2:", c2()) // independent state
}

// 2. Go 1.22+: every iteration of a for loop has its own variable.
func perIterationVariables() {
	fmt.Println("== 2. per-iteration loop variables (Go 1.22+) ==")
	var funcs []func() int
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() int { return i })
	}
	for _, f := range funcs {
		fmt.Print("  ", f())
	}
	fmt.Println()

	var ptrs []*int
	for i := range 3 { // range over int (Go 1.22)
		ptrs = append(ptrs, &i)
	}
	fmt.Println("  distinct addresses:", ptrs[0] != ptrs[1] && ptrs[1] != ptrs[2])
}

// 3. Simulating the pre-1.22 behaviour: one variable shared by all closures.
func sharedVariable() {
	fmt.Println("== 3. one shared variable (what Go <= 1.21 did) ==")
	var funcs []func() int
	i := 0
	for ; i < 3; i++ {
		funcs = append(funcs, func() int { return i })
	}
	for _, f := range funcs {
		fmt.Print("  ", f())
	}
	fmt.Println()
}

// 4. Per-iteration copy semantics in a 3-clause loop:
// each new iteration's i is initialised from the previous iteration's i
// BEFORE the post statement runs, so changes inside the body carry over.
func copySemantics() {
	fmt.Println("== 4. body changes carry over to the next iteration ==")
	var funcs []func() int
	for i := 0; i < 6; i++ {
		funcs = append(funcs, func() int { return i })
		i++ // modifies this iteration's variable; next iteration starts from it
	}
	for _, f := range funcs {
		fmt.Print("  ", f())
	}
	fmt.Println()
}

// 5. Goroutines capturing loop variables are now safe (still need sync!).
func goroutines() {
	fmt.Println("== 5. goroutines in a loop ==")
	var (
		mu   sync.Mutex
		wg   sync.WaitGroup
		seen []string
	)
	for _, name := range []string{"a", "b", "c"} {
		wg.Go(func() { // Go 1.25 WaitGroup.Go
			mu.Lock()
			seen = append(seen, name)
			mu.Unlock()
		})
	}
	wg.Wait()
	sort.Strings(seen)
	fmt.Println("  seen:", seen)
}

type Account struct{ Balance int }

func (a Account) Value() int    { return a.Balance }
func (a *Account) Pointer() int { return a.Balance }

// 6. Method values bind the receiver when the method value is evaluated.
func methodValues() {
	fmt.Println("== 6. method values bind the receiver early ==")
	acc := Account{Balance: 100}
	byValue := acc.Value     // copies acc NOW
	byPointer := acc.Pointer // takes &acc NOW
	acc.Balance = 999
	fmt.Println("  value receiver  :", byValue())
	fmt.Println("  pointer receiver:", byPointer())
}

// 7. defer: arguments are evaluated immediately, closure bodies later.
func deferEvaluation() {
	fmt.Println("== 7. defer argument vs closure ==")
	x := 1
	defer fmt.Println("  deferred arg    :", x)
	defer func() { fmt.Println("  deferred closure:", x) }()
	x = 2
}

func main() {
	captureByReference()
	perIterationVariables()
	sharedVariable()
	copySemantics()
	goroutines()
	methodValues()
	deferEvaluation()
}
