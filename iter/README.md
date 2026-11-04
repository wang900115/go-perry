# iter

Range-over-func iterators (`iter.Seq`, Go 1.23+) and the iterator helpers in `slices`, applied to a list of orders.

## What the code does

`main.go` defines an `Order` struct and five sample orders, then:

1. `filter(seq, keep)`: a generic function that wraps an `iter.Seq[V]` and yields only matching values. Used to print orders with `Amount > 50` (orders 1 and 5).
2. Filters again for `Status == "delivered"` (orders 1, 2, 5).
3. Writes a hand-rolled `func(yield func(Order) bool)` sequence for `Amount < 20` and materializes it with `slices.Collect` (prints `3: 12.00`).
4. `slices.SortedFunc(slices.Values(orders), cmp)` sorts by amount, descending (order 1, 5, 2, 4, 3).
5. `slices.Chunk(orders, 3)` prints the orders in groups of 3 and 2.

Actual output:

```
Order ID: 1, Customer: Alice, Amount: 75.50, Status: delivered
Order ID: 5, Customer: Eve, Amount: 71.22, Status: delivered
Order ID: 1, Customer: Alice, Amount: 75.50, Status: delivered
Order ID: 2, Customer: Bob, Amount: 45.90, Status: delivered
Order ID: 5, Customer: Eve, Amount: 71.22, Status: delivered
3: 12.00
OrderID: 1, Customer: Alice, Amount: 75.50, Status: delivered
OrderID: 5, Customer: Eve, Amount: 71.22, Status: delivered
OrderID: 2, Customer: Bob, Amount: 45.90, Status: delivered
OrderID: 4, Customer: Dan, Amount: 30.99, Status: canceled
OrderID: 3, Customer: Candy, Amount: 12.00, Status: pending
[{1 Alice 75.5 delivered} {2 Bob 45.9 delivered} {3 Candy 12 pending}]
[{4 Dan 30.99 canceled} {5 Eve 71.22 delivered}]
```

## Concepts you can learn

- **Push iterators**: an `iter.Seq[V]` is just `func(yield func(V) bool)`, and `for v := range seq` calls it.
- **Respecting `yield`'s return value**: when the consumer `break`s, `yield` returns `false` and the iterator must stop; both `filter` and the hand-written sequence do this.
- **Lazy vs eager** (the comment at the top of the file): `filter` does no work until ranged over and composes without building intermediate slices, while `slices.Collect` and `slices.SortedFunc` consume the whole sequence into a new slice.
- **`cmp.Compare` with swapped arguments** (`b, a`) for a descending sort.
- **Generic iterator adapters**: `filter[V any]` works for any element type.

## How to run

This directory has no `go.mod`; it belongs to the root module (which declares a Go version new enough for `iter`). From the repo root:

```bash
go run iter/main.go
```
