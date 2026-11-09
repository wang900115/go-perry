# 2026-11-09 · Error wrapping: %w, errors.Is / As / AsType, error trees

## Why learn this
- **Every layered Go program**: repository → service → handler. Each layer adds context ("service: get user: ...") but the top layer still needs to ask *"was this a not-found?"* to choose HTTP 404 vs 500, or *"was it a timeout?"* to decide whether to retry.
- **Standard library APIs depend on it**: `errors.Is(err, fs.ErrNotExist)`, `errors.Is(err, context.DeadlineExceeded)`, `errors.As(err, &*net.OpError)`, `errors.Is(err, sql.ErrNoRows)`.
- **Collecting several failures**: validation of many fields, closing several resources, fan-out work (the repo's `errgroup` returns only the first error; `errors.Join` keeps all of them).

This follows the recent standard-library series (`sync_*`, `os_root`). The example in `os_root` already used `errors.Is`; this lesson explains how it actually works, and uses the new Go 1.26 `errors.AsType`.

## Core concepts
1. **Wrapping**: `fmt.Errorf("ctx: %w", err)` returns an error whose `Unwrap() error` returns `err`. `%v` / `%s` only copy the text, so the chain is lost (https://pkg.go.dev/fmt#Errorf).
2. **`errors.Is(err, target)`**: is `target` anywhere in the chain? Use it instead of `err == target` (https://pkg.go.dev/errors#Is).
3. **`errors.As(err, &target)`**: find the first error in the chain whose *type* matches, and assign it. `target` must be a non-nil pointer to a type that implements `error` (or to an interface), otherwise it panics (https://pkg.go.dev/errors#As).
4. **`errors.AsType[E](err) (E, bool)`** (Go 1.26): the generic version. No pre-declared variable, the type is checked at compile time, and no reflection is used (https://pkg.go.dev/errors#AsType, https://go.dev/doc/go1.26).
5. **Custom hooks**: a type may define `Is(error) bool`, `As(any) bool`, `Unwrap() error` or `Unwrap() []error` to customise matching.
6. **Trees** (Go 1.20): `errors.Join(errs...)` and `fmt.Errorf` with several `%w` produce an error with `Unwrap() []error`. `Is`/`As` search it depth-first. `errors.Unwrap` does **not** look into it (returns nil). `errors.Join` drops nil arguments and returns nil if all are nil; its message joins with newlines (https://go.dev/doc/go1.20#wrapping).

### How it works (from `$GOROOT/src/errors/wrap.go`, go1.26)
`errors.Is` is a loop, not magic:
```go
for {
    if targetComparable && err == target { return true }      // 1. plain ==
    if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) { return true } // 2. custom Is
    switch x := err.(type) {
    case interface{ Unwrap() error }:   err = x.Unwrap()       // 3a. single chain: next node
    case interface{ Unwrap() []error }: for each child: recurse // 3b. tree: depth-first
    default: return false
    }
}
```
- Step 1 checks with `reflectlite` whether the target's dynamic type is comparable first, so comparing against an uncomparable error type does not panic.
- `errors.AsType` uses the same loop but step 1 is a plain type assertion `err.(E)`, which is why it needs no reflection. `errors.As` instead uses `reflectlite` to validate the target and to assign through the pointer.
- `fmt.Errorf` returns `*fmt.wrapError` (one `%w`) or `*fmt.wrapErrors` (several `%w`); `errors.Join` returns `*errors.joinError`. These are unexported; only their `Unwrap` methods matter. (Visible in the `%T` output of the example.)
- **Typed nil**: an interface value is a pair (dynamic type, data pointer). It is `nil` only if *both* are nil. Returning a nil `*QueryError` as `error` produces (`*QueryError`, nil), which is `!= nil`.

## Full example
See `main.go` in this directory. Verified with `go run .` and `go vet` on go1.26.3 linux/amd64; `GOOS=windows go vet` also passes. Requires Go 1.26 for `errors.AsType`.

Expected output: see [README.md](README.md). On Windows the OS message in Part 4 is different (likely `The system cannot find the file specified.`, not verified on a real Windows machine), but every `true`/`false` line is the same.

## Common pitfalls
1. **Using `%v` (or `err.Error()`) when you meant `%w`.** The message looks identical, but `errors.Is` silently returns false. Conversely, `%w` makes the wrapped error part of your API: callers may start depending on it, so only wrap what you are willing to expose.
2. **Returning a typed nil as `error`.** `var e *MyErr; return e` is never `== nil`. Return a literal `nil` on the success path, and declare variables as `error`, not as the concrete pointer type.
3. **Assuming a single chain.** `errors.Unwrap` returns nil on a joined error, so hand-written `for err != nil { err = errors.Unwrap(err) }` loops miss branches. Use `errors.Is` / `errors.As` / `errors.AsType`, which understand trees. Also: `errors.As(err, qe)` (forgetting `&`) panics at runtime; `AsType` turns that mistake into a compile error.

## Exercise
Write `func ValidateUser(u User) error` for `type User struct{ Name string; Age int; Email string }` that checks all three fields and returns **all** failures with `errors.Join`:
- Define a `FieldError{Field string; Err error}` type with `Error()` and `Unwrap()`, and sentinels `ErrRequired`, `ErrOutOfRange`.
- In `main`, call it with a user that fails every check. Then:
  - Print whether `errors.Is(err, ErrRequired)` is true.
  - Use `errors.AsType[*FieldError]` and explain why it finds only the **first** `FieldError`.
  - Write a helper `FieldErrors(err error) []*FieldError` that collects **every** `FieldError` in the tree (hint: walk `Unwrap() []error` yourself).
- Make sure a valid user returns exactly `nil` (check with `== nil`).
