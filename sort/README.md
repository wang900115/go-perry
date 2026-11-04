# sort

Sorting with the standard `sort` package: (1) normal sorts of strings, ints, and floats, (2) sorting a slice of structs, and (3) implementing `sort.Interface`.

## What the code does

- `01.go` - **Normal sort (string, int, float)**. Uses `sort.Sort` with `sort.StringSlice`, `sort.IntSlice`, and `sort.Float64Slice`, and wraps them in `sort.Reverse` for descending order. It also sorts the strings by length with `sort.Slice`, then with `sort.SliceStable` ascending and descending. The commented-out lines show the shortcuts `sort.Strings`, `sort.Ints`, and `sort.Float64s`.
- `02.go` - **Struct sort**. Sorts `[]User{Name, Age}` by age ascending with `sort.Slice`, then descending with `sort.SliceStable`.
- `03.go` - **Implement the sort interface**. `UsersByAge` implements `Len`, `Swap`, and `Less` (by age). It is sorted with `sort.Sort`, then again with `sort.Stable`.

Actual output of `01.go`:

```
Sorted string: [brown fox jumps quick]
Reversed: [quick jumps fox brown]
Sorted String by length: [fox quick jumps brown]
[Stable] Sorted String by length: [fox quick jumps brown]
[Stable] Reversd Sorted String by length: [quick jumps brown fox]
Sorted int: [14 19 25 56 67 78]
Reversed: [78 67 56 25 19 14]
Sorted float: [19.5 20.8 57.4 176.8]
Reversed: [176.8 57.4 20.8 19.5]
```

`02.go` prints users ordered `Joda 12 ... Kap 58`, then reversed. `03.go` prints `[{Amanda 16} {Rajeev 28} {Steven 28} {Monica 31} {John 56}]` twice.

## Concepts you can learn

- **`sort.Interface` (`Len`/`Less`/`Swap`)**: the generic contract that `sort.Sort` works on. Named slice types like `StringSlice` and your own `UsersByAge` plug into it.
- **`sort.Reverse`**: wraps an `Interface` and flips `Less`, so you get descending order without a second comparator.
- **`sort.Slice` vs `sort.SliceStable`, and `Sort` vs `Stable`**: the unstable versions may reorder equal elements. The stable versions keep their original relative order (`O(n log n)` comparisons, at some extra cost). This matters when you sort by one key after another.
- **Closures as comparators**: `sort.Slice(x, func(i, j int) bool {...})` avoids defining a new type for one-off orderings.

## How to run

Part of the root module (no `go.mod` here). Each numbered file has its own `func main`, so run them one at a time from the repo root:

```bash
go run sort/01.go
go run sort/02.go
go run sort/03.go
```

`go build ./sort` fails by design.

## Notes

- Among strings of equal length, the order produced by `sort.Slice` (unstable) is not guaranteed. The output above is what this Go version produced, not a promise.
- In `03.go`, Rajeev and Steven are both 28. After the unstable `sort.Sort` their order is not guaranteed, so the following `sort.Stable` call only preserves whatever order `sort.Sort` left them in. It does not restore the original input order.
- Since Go 1.21, the `slices` package (`slices.Sort`, `slices.SortFunc`, `slices.SortStableFunc`) is the generic alternative. It is not used here.
