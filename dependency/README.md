# dependency

Dependency injection through interfaces. A `Bakery` depends on an `Oven` interface and a list of `Ingredient` interfaces, not on concrete types.

## What the code does

`main.go`:

- Defines two interfaces, `Oven` (`Heat()`) and `Ingredient` (`Mix()`).
- Implements them with `GasOven`, `ElectricOven`, `Flour`, `Sugar`, and `Butter`.
- `Bakery.Bake()` calls `oven.Heat()`, then `Mix()` on each ingredient, and prints "Baking an asesome pastry!".

`main` builds two bakeries with different dependencies and prints:

```
Heating with Gas Oven!
Mixing flour!
Mixing sugar!
Baking an asesome pastry!
Heating with electrice oven!
Mixing flour!
Mixing butter!
Baking an asesome pastry!
```

## Concepts you can learn

- **Dependency injection via struct literal**: the caller decides which `Oven` and which ingredients to use (a `NewBakery(oven, ingredients...)` constructor would be the usual alternative). `Bakery` never creates them itself.
- **Programming to interfaces**: you can swap `GasOven` for `ElectricOven` (or a test fake) without changing `Bakery`.
- **Implicit interface satisfaction**: the concrete types never declare `implements`. Having the method is enough.

## How to run

`dependency/` has no `go.mod`, so it belongs to the root module. Run it from the repository root:

```bash
go run ./dependency
```

## Notes

- The output strings contain typos ("asesome", "electrice"). They come from the source as written.
