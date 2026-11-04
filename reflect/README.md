# reflect

Basics of the `reflect` package: inspecting types and values, reading and writing struct fields, calling methods by name, and creating values at runtime.

## What the code does

Each file is a standalone program with its own `func main`:

- `01.go`: `reflect.TypeOf`, `reflect.ValueOf` and `reflect.Zero` on an `int`. Prints `Type: int`, `Value: 0`, `Zero Value: 0`.
- `02.go`: reads `Name`/`Age` from a `Person` with `FieldByName`. It then gets a settable value with `reflect.ValueOf(&person).Elem()` and calls `SetInt(15)`. Prints `Name: Alan`, `Age: 30`, `Updated Age: 15`.
- `03.go`: looks up `Calculator.Add` with `MethodByName` and calls it with `[]reflect.Value{5, 3}`. It turns the result back into an `int` with `.Interface().(int)` and prints `Result: 8`.
- `04.go`: `reflect.New(reflect.TypeOf(0)).Elem()` creates a new addressable `int`, sets it to 42 and prints `New Integer Value: 42`.

## Concepts you can learn

- **`Type` vs `Value`**: `TypeOf` describes the static shape and `ValueOf` wraps a concrete value. Most reflection code moves between the two.
- **Settability**: `ValueOf(x)` wraps a copy and cannot be set. You need a pointer plus `.Elem()` to modify the original variable (`02.go`).
- **Dynamic invocation**: `MethodByName` + `Call` let code invoke methods chosen at runtime. This is the basis of RPC dispatchers and template engines. Results come back as `[]reflect.Value` and must be type-asserted.
- **Creating values**: `reflect.New` allocates a zero value of a type known only at runtime, which decoders use to build target objects.

## How to run

This directory has no `go.mod`; it belongs to the root module. From the repo root:

```bash
go run reflect/01.go
go run reflect/02.go
go run reflect/03.go
go run reflect/04.go
```

`go build ./reflect` fails by design, because every file declares `main`.

## Notes

- Code comments are in Chinese.
- `FieldByName`/`MethodByName` panic or return an invalid `Value` on a misspelled name. The examples do not check for this.
