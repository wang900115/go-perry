# struct_embedding

How struct embedding works in Go: an embedded field's fields and methods are
*promoted* to the outer type, but this is composition, not inheritance. Covers
promotion and the explicit path (`u.Base.Name`), method sets when embedding `T`
vs `*T`, why a shadowing method does not "override" calls made inside the
embedded type, the depth rule and ambiguous selectors, embedding an interface to
build a decorator, how embedding `sync.Mutex` leaks `Lock`/`Unlock` into your
public API, and how `encoding/json` flattens embedded structs.

```bash
go run .
```

Standard library only. Verified with go1.26.3 linux/amd64; `go vet .` and
`GOOS=windows go vet .` pass. The commented-out lines in `main.go` are real
compile errors (checked by uncommenting them).
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: field and method promotion ==
u.Name       = perry
u.Base.Name  = perry
u.Describe() = Base(id=1, name=perry)
after Rename = Perry

== Part 2: method sets ==
*ByValue  implements Renamer: *main.ByValue
ByPointer implements Renamer: main.ByPointer, name=b2
recovered: runtime error: invalid memory address or nil pointer dereference

== Part 3: shadowing, not virtual dispatch ==
d.Sound()        = woof
d.Animal.Sound() = ...
d.Speak()        = dog says ...
via interface    = dog says ...

== Part 4: depth rule and ambiguity ==
s.Logger.Level  = debug
s.Metrics.Close = Metrics.Close
Service has Close(): false
w.Level (depth 0) = warn

== Part 5: embedding an interface ==
counted=17 buffer="hello gopher\nbye\n"
recovered: runtime error: invalid memory address or nil pointer dereference

== Part 6: public API and JSON ==
Counter.n = 1 SafeCounter.n = 1
{"created_by":"perry","version":2,"title":"embedding","owner":{"created_by":"team","version":1}}
```
