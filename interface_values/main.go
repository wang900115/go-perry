// Interface values in Go: what an interface really stores.
//
// Run: go run .
// Allocations: go build -gcflags=-m . 2>&1 | grep "escapes to heap"
package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"unsafe"
)

type Shape interface {
	Area() float64
}

type Square struct{ Side float64 }

func (s Square) Area() float64 { return s.Side * s.Side }

type Circle struct{ R float64 }

// Pointer receiver: only *Circle has Area in its method set.
func (c *Circle) Area() float64 { return 3 * c.R * c.R }

// Compile-time checks: the build fails if a type stops satisfying Shape.
var (
	_ Shape = Square{}
	_ Shape = (*Circle)(nil)
	// _ Shape = Circle{} // compile error: Area has pointer receiver
)

// eface mirrors runtime.eface (the layout of an `any` value).
// This is for illustration only; never rely on it in real code.
type eface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

// 1. An interface value is two words: (type, data).
func twoWords() {
	fmt.Println("== 1. an interface is two words ==")
	var s Shape
	var a any
	fmt.Println("  sizeof(Shape) =", unsafe.Sizeof(s), " sizeof(any) =", unsafe.Sizeof(a))

	a = 42
	e := (*eface)(unsafe.Pointer(&a))
	fmt.Println("  any(42): type word set =", e.typ != nil, ", data word set =", e.data != nil)
	fmt.Println("  *(*int)(data) =", *(*int)(e.data))
}

// 2. Storing a value copies it. Later changes to the original are not seen.
func copiesValue() {
	fmt.Println("== 2. assignment copies the value ==")
	sq := Square{Side: 2}
	var s Shape = sq
	sq.Side = 10
	fmt.Println("  original area =", sq.Area(), " interface area =", s.Area())

	c := &Circle{R: 1}
	var s2 Shape = c // copies the pointer, not the Circle
	c.R = 2
	fmt.Println("  pointer stored: interface sees the change, area =", s2.Area())
}

// 3. Boxing can allocate. Small integers and zero-sized values do not.
var sink any

func boxing() {
	fmt.Println("== 3. boxing and allocation ==")
	// Values are computed at run time so the compiler cannot use a static copy.
	i := 0
	small := testing.AllocsPerRun(1000, func() { i = (i + 1) % 200; sink = i })  // 0..199
	big := testing.AllocsPerRun(1000, func() { i = 1000 + (i+1)%200; sink = i }) // >= 256
	words := strings.Fields("alpha beta gamma")
	str := testing.AllocsPerRun(1000, func() { i = (i + 1) % 3; sink = words[i] }) // string header
	p := &Square{Side: 1}
	ptr := testing.AllocsPerRun(1000, func() { sink = p })            // pointer fits the data word
	empty := testing.AllocsPerRun(1000, func() { sink = struct{}{} }) // zero-sized
	fmt.Printf("  small int (<256): %.0f allocs/op\n", small)
	fmt.Printf("  big int:          %.0f allocs/op\n", big)
	fmt.Printf("  string:           %.0f allocs/op\n", str)
	fmt.Printf("  pointer:          %.0f allocs/op\n", ptr)
	fmt.Printf("  struct{}{}:       %.0f allocs/op\n", empty)
}

// 4. Type assertions and type switches read the type word.
func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil interface"
	case int, int64:
		return fmt.Sprintf("integer %v (x is still any here: %T)", x, x)
	case fmt.Stringer:
		return "Stringer: " + x.String()
	case error:
		return "error: " + x.Error()
	default:
		return fmt.Sprintf("other %T", x)
	}
}

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

func assertions() {
	fmt.Println("== 4. type assertions and type switches ==")
	var s Shape = Square{Side: 3}

	sq, ok := s.(Square)
	fmt.Println("  s.(Square) ->", sq, ok)
	_, ok = s.(*Circle)
	fmt.Println("  s.(*Circle) ok =", ok)

	// Asserting to another interface: checks the method set at run time.
	_, ok = s.(fmt.Stringer)
	fmt.Println("  s.(fmt.Stringer) ok =", ok)

	for _, v := range []any{nil, 5, Celsius(21.5), errors.New("boom"), []int{1}} {
		fmt.Println("  describe:", describe(v))
	}

	defer func() { fmt.Println("  recovered:", recover()) }()
	_ = s.(*Circle) // single-value form panics on mismatch
}

// 5. Optional interfaces: io.Copy upgrades behavior if the value has more methods.
type plainReader struct{ r io.Reader }

func (p plainReader) Read(b []byte) (int, error) { return p.r.Read(b) }

func optionalInterfaces() {
	fmt.Println("== 5. optional interfaces ==")
	var r io.Reader = strings.NewReader("data")
	_, isWT := r.(io.WriterTo)
	fmt.Println("  *strings.Reader implements io.WriterTo:", isWT)
	_, isWT = io.Reader(plainReader{r}).(io.WriterTo)
	fmt.Println("  wrapper implements io.WriterTo:      ", isWT, "(wrapping hides methods)")
}

// 6. Comparing interfaces: equal if same dynamic type AND equal values.
// If the dynamic type is not comparable, == panics at run time.
func comparison() {
	fmt.Println("== 6. comparing interface values ==")
	var a, b any = 1, 1
	var c any = int64(1)
	fmt.Println("  any(1) == any(1):        ", a == b)
	fmt.Println("  any(1) == any(int64(1)): ", a == c)

	m := map[any]string{}
	m[1] = "int"
	m["1"] = "string"
	fmt.Println("  map[any] with 1 and \"1\": len =", len(m))

	defer func() { fmt.Println("  recovered:", recover()) }()
	var x, y any = []int{1}, []int{1}
	fmt.Println(x == y) // panics: slices are not comparable
}

// 7. nil interface vs interface holding a nil pointer.
func find(ok bool) Shape {
	var c *Circle
	if ok {
		c = &Circle{R: 1}
	}
	return c // always a non-nil interface: the type word is *Circle
}

func nilInterface() {
	fmt.Println("== 7. nil interface vs typed nil ==")
	var s Shape
	fmt.Println("  zero Shape == nil:", s == nil)
	r := find(false)
	fmt.Printf("  find(false) == nil: %v (dynamic type %T)\n", r == nil, r)
}

func main() {
	twoWords()
	copiesValue()
	boxing()
	assertions()
	optionalInterfaces()
	comparison()
	nilInterface()
}
