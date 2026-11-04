# decorator

The Decorator design pattern in Go, built by composition: each decorator holds a `Coffee` interface field and delegates to it. It uses the classic coffee-and-add-ons example.

## What the code does

`main.go` defines a `Coffee` interface (`Cost()`, `Description()`) and three implementations:

- `SimpleCoffee` costs 2.0 with the description "Simple Coffee".
- `Milk` wraps a `Coffee`, adds 0.5, and appends ", Milk".
- `Caramel` wraps a `Coffee`, adds 1, and appends ", Caramel".

`main` builds three combinations and prints:

```
Coffee with Milk:  Simple Coffee, Milk Cost:  2.5
Coffee with Caramel:  Simple Coffee, Caramel Cost:  3
Coffee with Milk and Caramel: Simple Coffee, Caramel, Milk Cost:  3.5
```

## Concepts you can learn

- **Decorator pattern**: you add behavior by wrapping an object that has the same interface, instead of subclassing. Wrappers can be stacked in any order and combined at runtime.
- **Composition over inheritance in Go**: each decorator holds a `Coffee` field and delegates to it. Go has no class inheritance, so this is the idiomatic way to do it.
- **Wrapping order matters**: `Milk{Caramel{SimpleCoffee}}` yields "Caramel, Milk", because the outermost wrapper runs last.

## How to run

`decorator/` has no `go.mod`, so it belongs to the root module. Run it from the repository root:

```bash
go run ./decorator
```
