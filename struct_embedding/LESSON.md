# 2026-11-20 · Struct embedding: promotion, method sets and shadowing

## Why learn this
- Embedding is everywhere in real Go code: `bufio.ReadWriter` embeds `*Reader` and `*Writer`, HTTP middleware wraps `http.ResponseWriter` by embedding it, test fakes embed an interface and implement only one method, and many structs embed `sync.Mutex`.
- It looks like inheritance, but it is not. Developers coming from Java/C# expect "override" and virtual dispatch and get silent bugs (Part 3).
- It connects earlier topics: method sets decide interface satisfaction (`error_wrapping`, `io_reader_writer`), and `encoding/json` treats embedded fields specially (`json_encoding`).

## Core concepts
1. **An embedded field is a field declared with a type but no name** (`T` or `*T`, where `T` is not a pointer or interface pointer). Its field name is the type name: `u.Base` (https://go.dev/ref/spec#Struct_types).
2. **Promotion:** fields and methods of the embedded field can be used as if they belonged to the outer struct: `u.Name` is shorthand for `u.Base.Name`. Promoted fields cannot be used as names in a composite literal: write `User{Base: Base{Name: "x"}}`, not `User{Name: "x"}`.
3. **Method sets** (https://go.dev/ref/spec#Method_sets and Struct types):
   - Embedding `T`: `S` gets the methods of `T` with value receivers; `*S` gets the methods of `T` and `*T`.
   - Embedding `*T`: both `S` and `*S` get the methods of `T` and `*T`.
   - So `ByValue{}` does not implement an interface whose method has a pointer receiver, but `&ByValue{}` and `ByPointer{}` do.
4. **Selectors and depth** (https://go.dev/ref/spec#Selectors): `x.f` picks the `f` at the *shallowest* depth. If there are two `f` at the same shallowest depth, `x.f` is an "ambiguous selector" compile error — but only if you use it. Ambiguous methods are also not in the method set, so the type silently does not satisfy an interface that needs them.
5. **Shadowing, not overriding:** if the outer type declares a method with the same name, it hides the promoted one. But a promoted method like `Animal.Speak` still has an `Animal` receiver; calls inside it go to `Animal.Sound`. There is no virtual dispatch on structs. If you need "template method" behaviour, pass an interface or a function in.
6. **Embedding an interface** in a struct promotes all interface methods; the struct implements the interface automatically, and you replace only the methods you care about. If the embedded interface value is nil, calling an un-replaced method panics.
7. **Embedding is part of your API.** Exported methods of an embedded type become exported methods of your type (`c.Lock()` on a struct that embeds `sync.Mutex`). Prefer a named, unexported field (`mu sync.Mutex`) unless you really want to expose them.
8. **`encoding/json` flattens anonymous struct fields** into the outer object, following Go's visibility rules; a named field (`Owner Meta`) becomes a nested object (https://pkg.go.dev/encoding/json#Marshal). If two embedded structs produce the same JSON key at the same depth and neither is tagged, the field is silently omitted (see the "multiple fields" rule in that doc).

### How it works (compiler/runtime)
- Promotion is resolved **at compile time**. `u.Name` compiles to a fixed offset load: `u.Base` is just the first bytes of `User` (embedding a value type means the struct is laid out inline, no pointer). `u.Describe()` compiles to `Base.Describe(u.Base)`.
- For an embedded `*T`, `s.Name` is a pointer dereference, which is why a nil embedded pointer panics with a nil pointer dereference at the point of use (Part 2).
- When a type with promoted methods is converted to an interface, the compiler generates small **wrapper methods** (e.g. `main.ByPointer.Rename` which loads the `*Base` field and calls `(*Base).Rename`). The itab points at these wrappers. You can see them in `go tool nm` output or with `go build -gcflags=-S`. This is an implementation detail; I am not certain of the exact naming across versions.
- None of this involves the runtime type system at call time: there is no vtable lookup through the embedding chain.

## Full example
See `main.go`. Verified with `go run .` (go1.26.3 linux/amd64), `go vet .` and `GOOS=windows go vet .`. The commented-out lines were uncommented one by one to confirm the compile errors. Expected output is in [README.md](README.md).

Parts:
1. Field and method promotion, explicit path, calling a pointer method on an addressable value.
2. Method sets: embed `T` vs `*T`, interface satisfaction, nil embedded pointer panic.
3. A "shadowing" `Dog.Sound` is not used by the promoted `Animal.Speak`.
4. Ambiguous selectors at the same depth, ambiguous methods dropped from the method set, depth 0 wins.
5. `CountingWriter` embeds `io.Writer` (decorator) and the nil-interface panic.
6. Embedding `sync.Mutex` exposes `Lock`/`Unlock`; JSON flattening of an embedded struct vs a named field.

## Common pitfalls
1. **Expecting overriding.** Writing `Dog.Sound` and expecting `Animal.Speak` to call it. It never will. Design with interfaces or function fields when the base type must call back into the outer type.
2. **Nil embedded pointers and interfaces.** `struct{ *Base }` or `struct{ io.Writer }` compile with a nil zero value, and the first promoted call panics. Construct with a `New...` function, or embed by value.
3. **Leaking or copying embedded state.** Embedding `sync.Mutex` adds `Lock`/`Unlock` to your public API, and copying the outer struct copies the mutex (`go vet` reports this as `copylocks`). Also, adding a method to an embedded type later can create a new ambiguity or silently change which method is promoted — embedding couples you to that type's API.

## Exercise (no answer given)
Write an HTTP middleware that records the status code, using only the standard library and `net/http/httptest` (no real server needed):
- Define `type statusRecorder struct { http.ResponseWriter; status int }` and override `WriteHeader` to store the status before calling the embedded one.
- Write `func logStatus(next http.Handler) http.Handler` that prints `METHOD PATH -> STATUS`.
- Test it with `httptest.NewRecorder()` and a handler that (a) calls `w.WriteHeader(404)`, (b) only calls `w.Write` without `WriteHeader`. What status does your recorder report in case (b), and why? Fix it.
- Bonus: your wrapper hides optional interfaces such as `http.Flusher`. Check with a type assertion `w.(http.Flusher)` inside the handler and explain the result. Then read about `http.ResponseController` (https://pkg.go.dev/net/http#ResponseController) and the `Unwrap() http.ResponseWriter` convention, and add `Unwrap` to your recorder.

## Suggested directory
`struct_embedding`
