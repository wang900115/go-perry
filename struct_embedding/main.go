// Command struct_embedding shows how struct embedding works in Go:
// field and method promotion, method sets, shadowing (not overriding),
// ambiguous selectors, embedding interfaces, and JSON flattening.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

func main() {
	part1()
	part2()
	part3()
	part4()
	part5()
	part6()
}

// ---------- Part 1: promotion ----------

type Base struct {
	ID   int
	Name string
}

func (b Base) Describe() string { return fmt.Sprintf("Base(id=%d, name=%s)", b.ID, b.Name) }
func (b *Base) Rename(n string) { b.Name = n }

type User struct {
	Base  // embedded field: its name is the type name "Base"
	Email string
}

func part1() {
	fmt.Println("== Part 1: field and method promotion ==")
	u := User{Base: Base{ID: 1, Name: "perry"}, Email: "p@example.com"}
	fmt.Println("u.Name       =", u.Name)      // promoted field
	fmt.Println("u.Base.Name  =", u.Base.Name) // explicit path, same field
	fmt.Println("u.Describe() =", u.Describe())
	u.Rename("Perry") // (&u.Base).Rename: u is addressable, so this is allowed
	fmt.Println("after Rename =", u.Name)
	fmt.Println()
}

// ---------- Part 2: method sets and interfaces ----------

type Renamer interface{ Rename(string) }

type ByValue struct{ Base }    // embeds Base (value)
type ByPointer struct{ *Base } // embeds *Base (pointer)

func part2() {
	fmt.Println("== Part 2: method sets ==")
	var r Renamer

	// ByValue's method set contains only Describe (value receiver).
	// Rename (pointer receiver) is only in *ByValue's method set.
	// r = ByValue{}      // compile error: ByValue does not implement Renamer
	r = &ByValue{Base{Name: "a"}}
	r.Rename("a2")
	fmt.Printf("*ByValue  implements Renamer: %T\n", r)

	// ByPointer embeds *Base, so both Describe and Rename are promoted
	// into the method set of the value type ByPointer.
	p := ByPointer{&Base{Name: "b"}}
	r = p
	r.Rename("b2")
	fmt.Printf("ByPointer implements Renamer: %T, name=%s\n", r, p.Name)

	// A nil embedded pointer compiles fine but panics on use.
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		var empty ByPointer
		fmt.Println(empty.Describe())
	}()
	fmt.Println()
}

// ---------- Part 3: shadowing is NOT overriding ----------

type Animal struct{ Kind string }

func (a Animal) Sound() string { return "..." }
func (a Animal) Speak() string { return a.Kind + " says " + a.Sound() }

type Dog struct{ Animal }

func (d Dog) Sound() string { return "woof" } // shadows Animal.Sound

type Speaker interface{ Speak() string }

func part3() {
	fmt.Println("== Part 3: shadowing, not virtual dispatch ==")
	d := Dog{Animal{Kind: "dog"}}
	fmt.Println("d.Sound()        =", d.Sound())        // Dog.Sound
	fmt.Println("d.Animal.Sound() =", d.Animal.Sound()) // still reachable
	// Speak is promoted from Animal. Inside it, the receiver is an Animal,
	// so a.Sound() calls Animal.Sound, never Dog.Sound.
	fmt.Println("d.Speak()        =", d.Speak())
	var s Speaker = d
	fmt.Println("via interface    =", s.Speak())
	fmt.Println()
}

// ---------- Part 4: depth and ambiguous selectors ----------

type Logger struct{ Level string }
type Metrics struct{ Level string }

func (Logger) Close() string  { return "Logger.Close" }
func (Metrics) Close() string { return "Metrics.Close" }

type Service struct {
	Logger
	Metrics
	Name string
}

type Wrapper struct {
	Service
	Level string // depth 0 wins over the two depth-1 Level fields
}

func part4() {
	fmt.Println("== Part 4: depth rule and ambiguity ==")
	s := Service{Logger{"debug"}, Metrics{"info"}, "api"}
	// fmt.Println(s.Level)   // compile error: ambiguous selector s.Level
	// s.Close()              // compile error: ambiguous selector s.Close
	fmt.Println("s.Logger.Level  =", s.Logger.Level)
	fmt.Println("s.Metrics.Close =", s.Metrics.Close())
	// Service does not have a Close method at all (ambiguous ones are
	// dropped), so it does not satisfy interface{ Close() string }.
	_, ok := any(s).(interface{ Close() string })
	fmt.Println("Service has Close():", ok)

	w := Wrapper{Service: s, Level: "warn"}
	fmt.Println("w.Level (depth 0) =", w.Level)
	fmt.Println()
}

// ---------- Part 5: embedding an interface (decorator) ----------

// CountingWriter embeds io.Writer: every io.Writer method is promoted,
// and we override only Write to add behaviour.
type CountingWriter struct {
	io.Writer
	N int
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.Writer.Write(p) // call the embedded value explicitly
	c.N += n
	return n, err
}

func part5() {
	fmt.Println("== Part 5: embedding an interface ==")
	var buf bytes.Buffer
	cw := &CountingWriter{Writer: &buf}
	fmt.Fprintf(cw, "hello %s\n", "gopher")
	io.WriteString(cw, "bye\n")
	fmt.Printf("counted=%d buffer=%q\n", cw.N, buf.String())

	// Embedding an interface you never set leaves it nil: calling it panics.
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		bad := &CountingWriter{}
		bad.Write([]byte("x"))
	}()
	fmt.Println()
}

// ---------- Part 6: exported API leaks, JSON flattening ----------

// Embedding sync.Mutex makes Lock/Unlock part of Counter's public API.
type Counter struct {
	sync.Mutex
	n int
}

// SafeCounter keeps the mutex as a named, unexported field instead.
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

func (c *SafeCounter) Inc() { c.mu.Lock(); c.n++; c.mu.Unlock() }

type Meta struct {
	CreatedBy string `json:"created_by"`
	Version   int    `json:"version"`
}

type Doc struct {
	Meta         // anonymous: fields are flattened into Doc's JSON
	Title string `json:"title"`
	Owner Meta   `json:"owner"` // named: nested object
}

func part6() {
	fmt.Println("== Part 6: public API and JSON ==")
	var c Counter
	c.Lock() // anyone holding a Counter can lock it from outside
	c.n++
	c.Unlock()
	var sc SafeCounter
	sc.Inc()
	fmt.Println("Counter.n =", c.n, "SafeCounter.n =", sc.n)

	d := Doc{Meta: Meta{"perry", 2}, Title: "embedding", Owner: Meta{"team", 1}}
	out, _ := json.Marshal(d)
	fmt.Println(string(out))
}
