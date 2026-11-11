// Command defer_panic_recover demonstrates how defer, panic and recover work
// together: LIFO order, argument evaluation time, turning a panic into an
// error, the "recover must be called directly by a deferred function" rule,
// panic(nil) since Go 1.21, and recovering inside goroutines.
package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

// ---------- Part 1: defer order and argument evaluation ----------

func deferOrder() {
	x := 1
	// Arguments are evaluated NOW (x == 1), the call runs later.
	defer fmt.Println("deferred with arg x =", x)
	// A closure captures the variable, so it sees the final value.
	defer func() { fmt.Println("deferred closure sees x =", x) }()
	for i := range 3 {
		defer fmt.Println("loop defer", i) // runs 2, 1, 0 (LIFO)
	}
	x = 42
	fmt.Println("function body done, x =", x)
}

// A deferred closure can change a named result after "return" set it.
func double() (n int) {
	defer func() { n *= 2 }()
	return 21 // n = 21, then the deferred func runs, then the caller sees 42
}

// ---------- Part 2: turn a panic into an error ----------

// ErrPanic is the sentinel used to wrap recovered panics.
var ErrPanic = errors.New("recovered panic")

// safeCall runs f and converts any panic into an error.
// It works only because err is a NAMED result the deferred func can set.
func safeCall(name string, f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// If the panic value is already an error, keep it in the chain.
			if e, ok := r.(error); ok {
				err = fmt.Errorf("%s: %w: %w", name, ErrPanic, e)
				return
			}
			err = fmt.Errorf("%s: %w: %v", name, ErrPanic, r)
		}
	}()
	f()
	return nil
}

// ---------- Part 3: where recover works and where it does not ----------

// helperRecover is called BY a deferred function, not AS one,
// so recover returns nil here and the panic keeps going.
func helperRecover() {
	fmt.Println("  helper: recover() =", recover())
}

func recoverInHelper() {
	defer func() { helperRecover() }()
	panic("not caught by helper")
}

// ---------- Part 5: goroutines ----------

// goSafe starts f in a goroutine and reports a panic as an error.
// A panic in a goroutine can only be recovered inside that goroutine.
func goSafe(wg *sync.WaitGroup, errs chan<- error, id int, f func()) {
	wg.Go(func() { // sync.WaitGroup.Go: Go 1.25+
		if err := safeCall(fmt.Sprintf("worker %d", id), f); err != nil {
			errs <- err
		}
	})
}

func main() {
	fmt.Println("== Part 1: defer order and evaluation ==")
	deferOrder()
	fmt.Println("double() =", double())

	fmt.Println("== Part 2: panic -> error ==")
	err := safeCall("string panic", func() { panic("boom") })
	fmt.Println(err)
	fmt.Println("  errors.Is(err, ErrPanic):", errors.Is(err, ErrPanic))

	err = safeCall("index", func() {
		s := []int{1, 2, 3}
		i := 5
		_ = s[i]
	})
	fmt.Println(err)
	var re runtime.Error
	fmt.Println("  is runtime.Error:", errors.As(err, &re))

	err = safeCall("nil map", func() {
		var m map[string]int
		m["a"] = 1
	})
	fmt.Println(err)

	fmt.Println("== Part 3: recover must be called directly by a deferred func ==")
	fmt.Println(safeCall("helper", recoverInHelper))

	// runtime.Goexit runs deferred calls, but it is not a panic:
	// recover returns nil, and the goroutine still exits.
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() { fmt.Println("  Goexit: recover() =", recover()) }()
		runtime.Goexit()
	}()
	<-done

	fmt.Println("== Part 4: panic(nil) since Go 1.21 ==")
	err = safeCall("nil", func() { panic(nil) })
	fmt.Println(err)
	var pne *runtime.PanicNilError
	fmt.Println("  is *runtime.PanicNilError:", errors.As(err, &pne))

	fmt.Println("== Part 5: recovering inside goroutines ==")
	var wg sync.WaitGroup
	errs := make(chan error, 3)
	for id := range 3 {
		goSafe(&wg, errs, id, func() {
			if id == 1 {
				panic(fmt.Sprintf("worker %d failed", id))
			}
		})
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		fmt.Println(e)
	}
	fmt.Println("main still running: ok")
}
