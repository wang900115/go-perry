// Interface values in Go: what an interface value really is (a type word plus
// a data word), how boxing allocates, type assertions and type switches,
// comparing interfaces, and the typed-nil trap.
package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"unsafe"
)

// eface mirrors the runtime layout of an empty interface (runtime.eface).
// This is an implementation detail used only for teaching; never rely on it.
type eface struct {
	typ  unsafe.Pointer // *runtime._type: the dynamic type
	data unsafe.Pointer // pointer to the value (or the value itself for pointer-shaped types)
}

func words(v any) eface { return *(*eface)(unsafe.Pointer(&v)) }

type Shape interface {
	Area() float64
}

type Rect struct{ W, H float64 }

func (r Rect) Area() float64 { return r.W * r.H }

type Circle struct{ R float64 }

func (c *Circle) Area() float64 { return 3 * c.R * c.R } // pointer receiver, pi ~ 3 to keep output short

// Compile-time checks: the build fails if the type stops satisfying Shape.
var (
	_ Shape = Rect{}
	_ Shape = (*Circle)(nil)
	// _ Shape = Circle{} // does not compile: Area has a pointer receiver
)

type MyErr struct{}

func (*MyErr) Error() string { return "my error" }

func mayFail(fail bool) error {
	var p *MyErr // typed nil pointer
	if fail {
		p = &MyErr{}
	}
	return p // BUG: a non-nil interface holding a nil *MyErr when fail == false
}

func mayFailFixed(fail bool) error {
	if fail {
		return &MyErr{}
	}
	return nil // a truly nil interface: both words are nil
}

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil interface"
	case int, int64:
		return fmt.Sprintf("integer %v (x is still any here: %T)", x, x)
	case string:
		return fmt.Sprintf("string of len %d", len(x))
	case Shape:
		return fmt.Sprintf("Shape with area %.0f", x.Area())
	case error:
		return "error: " + x.Error()
	default:
		return fmt.Sprintf("other type %T", x)
	}
}

func main() {
	fmt.Println("== 1. An interface value is two words ==")
	fmt.Println("sizeof(any)   =", unsafe.Sizeof(any(nil)))
	fmt.Println("sizeof(Shape) =", unsafe.Sizeof(Shape(nil)))
	var empty any
	w := words(empty)
	fmt.Println("nil any: typ == nil:", w.typ == nil, "data == nil:", w.data == nil)
	a, b := words(Rect{1, 2}), words(Rect{3, 4})
	fmt.Println("same dynamic type -> same type word:", a.typ == b.typ)
	fmt.Println("different values -> different data words:", a.data != b.data)
	c := &Circle{R: 1}
	fmt.Println("pointer stored directly in data word:", words(c).data == unsafe.Pointer(c))

	fmt.Println("\n== 2. Boxing copies the value ==")
	r := Rect{W: 2, H: 3}
	var s Shape = r // copy of r is stored in the interface
	r.W = 100
	fmt.Println("r.Area() =", r.Area(), " s.Area() =", s.Area())

	fmt.Println("\n== 3. Boxing may allocate ==")
	var sink any
	big, small := 1000, 0
	// The values change on every run so the compiler cannot box a constant at build time.
	fmt.Println("allocs any(int >= 256):", testing.AllocsPerRun(100, func() { big++; sink = big }))
	fmt.Println("allocs any(int < 256): ", testing.AllocsPerRun(100, func() { small = (small + 1) % 200; sink = small }))
	fmt.Println("allocs any(Rect):      ", testing.AllocsPerRun(100, func() { sink = r }))
	fmt.Println("allocs any(*Circle):   ", testing.AllocsPerRun(100, func() { sink = c }))
	_ = sink

	fmt.Println("\n== 4. Type assertions ==")
	var x any = "hello"
	str, ok := x.(string)
	fmt.Printf("x.(string) -> %q, %v\n", str, ok)
	num, ok := x.(int)
	fmt.Printf("x.(int)    -> %d, %v (zero value, no panic)\n", num, ok)
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		_ = x.(int) // single-result form panics on mismatch
	}()
	// Interface-to-interface assertion: checked at run time via the method set.
	var rd io.Reader = strings.NewReader("abc")
	_, isWriterTo := rd.(io.WriterTo)
	_, isWriter := rd.(io.Writer)
	fmt.Println("strings.Reader is io.WriterTo:", isWriterTo, " io.Writer:", isWriter)

	fmt.Println("\n== 5. Type switch ==")
	for _, v := range []any{nil, 42, int64(5), "go", Rect{2, 2}, &Circle{1}, Circle{1}, errors.New("boom"), 3.14} {
		fmt.Println(describe(v))
	}

	fmt.Println("\n== 6. Comparing interfaces ==")
	var i1, i2 any = 1, 1
	var i3 any = int64(1)
	fmt.Println("any(1) == any(1):       ", i1 == i2)
	fmt.Println("any(1) == any(int64(1)):", i1 == i3, "(different dynamic types)")
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		var s1, s2 any = []int{1}, []int{1}
		fmt.Println(s1 == s2) // compiles, panics at run time
	}()
	m := map[any]string{}
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		m[[]int{1}] = "slice key" // same panic: map keys are hashed via the dynamic type
	}()

	fmt.Println("\n== 7. The typed-nil trap ==")
	err := mayFail(false)
	fmt.Println("err == nil:", err == nil, " dynamic type:", fmt.Sprintf("%T", err))
	ew := words(err)
	fmt.Println("type word nil:", ew.typ == nil, " data word nil:", ew.data == nil)
	fmt.Println("fixed: err == nil:", mayFailFixed(false) == nil)
}
