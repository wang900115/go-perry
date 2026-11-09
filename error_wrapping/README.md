# error_wrapping

Go error wrapping end to end: `%w` vs `%v`, `errors.Is` / `errors.As`, the Go 1.26 generic `errors.AsType`, custom `Is` / `Unwrap` methods, multi-error trees (`errors.Join`, multiple `%w` in one `fmt.Errorf`), and the typed-nil interface pitfall.

## Run

```
go run .
```

Requires Go 1.26 (`errors.AsType`). Everything else works on Go 1.20+.

## What the code does

- **Part 1 – `%w` vs `%v`**: `==` fails on a wrapped sentinel, `errors.Is` succeeds; formatting with `%v` drops the chain.
- **Part 2 – `errors.As` vs `errors.AsType`**: extract a `*QueryError` from the chain, with and without a pre-declared target variable.
- **Part 3 – custom `Is`**: `HTTPError{404}` matches `ErrNotFound` without wrapping it.
- **Part 4 – trees**: `errors.Join` and `fmt.Errorf("...%w...%w")` build a tree; `errors.Is` searches it depth-first; `errors.Unwrap` returns nil for multi-errors.
- **Part 5 – typed nil**: returning a nil `*QueryError` as `error` gives a non-nil interface.

Expected output (observed on Linux with go1.26.3; `GOOS=windows go vet` also passes). On Windows the OS message in Part 4 differs (e.g. `The system cannot find the file specified.`), but every `true`/`false` is the same.

```
== 1. %w keeps the chain, %v breaks it ==
err: service: get user: query users id=42: not found
err == ErrNotFound         : false
errors.Is(err, ErrNotFound): true
lossy errors.Is            : false
errors.Unwrap(err)         : query users id=42: not found

== 2. errors.As vs errors.AsType (Go 1.26) ==
As     -> table=users id=42
AsType -> table=users id=42
AsType[*fs.PathError] found: false

== 3. custom Is method ==
404 is ErrNotFound: true
500 is ErrNotFound: false

== 4. multi-error trees ==
Join message (newline-separated):
service: get user: query users id=42: not found
cache: permission denied
Is ErrNotFound  : true
Is ErrPermission: true
errors.Unwrap(joined) == nil: true
errors.Join(nil, nil) == nil: true
multi-wrap message: load config: stat definitely-missing-file.txt: no such file or directory; fallback: fetch profile: http status 404
Is fs.ErrNotExist: true
Is ErrNotFound   : true
tree:
  - *fmt.wrapErrors: "load config: stat definitely-missing-file.txt: no such file or directory; fallback: fetch profile: http status 404"
    - *fs.PathError: "stat definitely-missing-file.txt: no such file or directory"
      - syscall.Errno: "no such file or directory"
    - *fmt.wrapError: "fetch profile: http status 404"
      - main.HTTPError: "http status 404"

== 5. typed-nil pitfall ==
validateBad(true)  != nil: true (dynamic type *main.QueryError)
validateGood(true) != nil: false
```

See [LESSON.md](LESSON.md) for the full lesson.
