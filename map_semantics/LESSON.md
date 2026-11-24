# Map Semantics

> Go version used: `go1.26.3`. The module's `go` line is `go 1.26.3`.
> Features used: `maps`/`clear` (Go 1.21), `maps.Keys` returning an iterator and `slices.Sorted` (Go 1.23), `range` over an int (Go 1.22).

## Why learn this

Maps are in almost every Go program: caches, counters, indexes by ID, sets, JSON objects (`map[string]any`),
HTTP headers (`http.Header` is a `map[string][]string`). Most map bugs are not about the hash table
algorithm (see the `swisstable` topic for that) but about the **language rules**:

- A struct field holding a `nil` map panics on the first write.
- A test passes locally and fails in CI because it depended on iteration order.
- "Why can't I write `m[k].Count++`?"
- A service crashes with `fatal error: concurrent map writes` and the `recover()` in the middleware does nothing.

## Core concepts

### 1. A map value is a pointer

`make(map[K]V)` allocates a runtime structure and returns a pointer to it. Since Go 1.24 that structure is a
Swiss-table map (`internal/runtime/maps`); before that it was `runtime.hmap`. Either way, the map value you hold
is a single pointer, so:

- `b := a` copies the pointer: both names see the same entries.
- Passing a map to a function lets the function modify the caller's map (no `*map` needed).
- `maps.Clone(m)` makes a new map with the same entries. It is **shallow**: if `V` is a pointer or slice,
  both maps share what it points to.

A `nil` map is a nil pointer. The runtime treats reads on it as reads on an empty map, so `len`, lookup and `range`
work; a write has nowhere to go and panics with `assignment to entry in nil map`.
([spec: Map types](https://go.dev/ref/spec#Map_types))

### 2. Comma-ok lookups

`m[k]` returns the zero value of `V` for a missing key. When the zero value is also a valid value
(`stock["apple"] == 0`), use `v, ok := m[k]`. ([spec: Index expressions](https://go.dev/ref/spec#Index_expressions))

### 3. Iteration order is unspecified, and the runtime randomizes it

The spec says the iteration order over maps is not specified and is not guaranteed to be the same from one
iteration to the next. ([spec: For statements with range clause](https://go.dev/ref/spec#For_range))
The runtime goes further and picks a random starting point for each `range`, precisely so that code cannot
accidentally depend on an order. For a deterministic order, sort the keys:

```go
for _, k := range slices.Sorted(maps.Keys(m)) { ... }
```

`maps.Keys` returns an `iter.Seq[K]` (Go 1.23), and `slices.Sorted` collects and sorts it.
([pkg maps](https://pkg.go.dev/maps), [pkg slices](https://pkg.go.dev/slices#Sorted))
Note that `fmt.Println(m)` prints maps with sorted keys (since Go 1.12), which can hide ordering bugs while debugging.
([pkg fmt](https://pkg.go.dev/fmt))

### 4. Map elements are not addressable

`&m[k]` and `m[k].X = 1` are compile errors (`cannot assign to struct field m["a"].X in map`). The reason: when the
map grows, the runtime moves entries to new memory, so a pointer to an element would dangle. Options:

- Copy out, modify, store back: `p := m[k]; p.X = 1; m[k] = p`.
- Store pointers: `map[string]*Point`, then `m[k].X = 1` modifies the pointed-to struct.

`m[k]++`, `m[k] += n` and `m[k] = append(m[k], v)` are fine: they are assignments to the map entry, not through
an element address. ([spec: Address operators](https://go.dev/ref/spec#Address_operators))

### 5. Mutating a map during `range`

From the spec: if an entry that has not yet been reached is removed during iteration, it will not be produced;
if an entry is created during iteration, it may be produced or skipped. So deleting while ranging is safe;
adding while ranging gives an unpredictable result. ([spec: For range](https://go.dev/ref/spec#For_range))

### 6. Keys must be comparable, and NaN is special

Key types must support `==`. Slices, maps and functions cannot be keys. With an interface key type (`map[any]V`),
the check moves to run time: inserting a `[]int` panics with `hash of unhashable type []int`.

Floating-point `NaN != NaN`, so every `m[NaN] = v` creates a **new** entry that can never be found or deleted.
`clear(m)` (Go 1.21) removes all entries, including NaN keys. ([spec: Clear](https://go.dev/ref/spec#Clear))

### 7. Concurrent access

Maps are not safe for concurrent use when at least one goroutine writes.
([Go FAQ: Why are map operations not defined to be atomic?](https://go.dev/doc/faq#atomic_maps))
The runtime has a cheap, best-effort detector: a "writing" flag is set during writes; if another read or write sees
it, the runtime calls `fatal("concurrent map writes")`. This is a **fatal error**, not a panic: deferred functions do
not run and `recover` cannot catch it. Because detection is best-effort, a racy program might not crash every time;
use `go test -race` / `go run -race` to find the race reliably. Protect maps with `sync.Mutex`/`sync.RWMutex`, or use
`sync.Map` for its specific use cases (see the `sync_map` topic).

## Full example

`main.go` in this directory. Run it with `go run .`.

Expected output (it is deterministic even though map iteration is not, because the program only prints sorted keys,
`fmt`-formatted maps, or booleans):

```
== 1. nil map: reads are fine, writes panic
len: 0 lookup: 0 nil? true
write err: assignment to entry in nil map

== 2. a map value is a pointer to the runtime map
a: map[x:1 y:2] b: map[x:1 y:2]
after addEntry(a): map[fromFunc:3 x:1 y:2]
a after editing clone: map[fromFunc:3 x:1 y:2] clone: map[fromFunc:3 x:100 y:2]

== 3. comma-ok: zero value vs missing key
apple=0 ok=true | banana=0 ok=false

== 4. iteration order is unspecified (and deliberately randomized)
distinct orders in 50 runs > 1: true
sorted keys: [a b c d e f g h]
fmt prints maps sorted: map[1:a 2:b 3:c]

== 5. map elements are not addressable
value map after store-back: map[p:{10 2}]
pointer map: {10 2}
word counts: map[fast:1 fun:1 go:3 is:2]

== 6. deleting during range is safe; adding is unpredictable
odd keys left: [1 3 5]
visited between 1 and 101 (exact count varies): true

== 7. NaN keys: you can insert, but never find or delete them
len: 2 found: false
len after clear: 0

== 8. set idiom and interface keys
set size: 2 has go: true
unhashable interface key err: runtime error: hash of unhashable type []int

== 9. concurrent writes
concurrent map writes are a fatal error, not a panic: recover() cannot catch it.
use sync.Mutex/RWMutex or sync.Map (see the sync_map topic).
```

"distinct orders in 50 runs > 1: true" is overwhelmingly likely but, strictly speaking, not guaranteed.

## Common pitfalls

1. **Forgetting to initialize a map field.** `type Cache struct{ m map[string]int }` followed by `c.m["k"] = 1`
   panics. Initialize in a constructor (`m: make(map[string]int)`) or lazily before the first write.
2. **Depending on iteration order.** Building a string, a hash, a SQL statement or test expectations by ranging over
   a map gives different results on different runs. Sort the keys first.
3. **Expecting `recover` to save you from a data race on a map.** `fatal error: concurrent map writes` kills the whole
   process. Synchronize access instead.

Bonus: deleting entries does not give memory back to the OS in the way you might expect. As far as I know, the
runtime does not shrink a map's tables after `delete` (this was true for the old `hmap`, and I believe it is still
true for the Swiss-table implementation, but I have not verified it in the source). If a map once held millions of
entries, copy the survivors into a new map to release the memory.

## Exercise

Write a function `func GroupByLength(words []string) map[int][]string` and a `main` that:

1. Groups words by length (`"go"`, `"is"` under 2; `"rust"`, `"java"` under 4; ...).
2. Prints the groups in ascending order of length, and the words within each group in alphabetical order,
   with exactly the same output on every run.
3. Adds a `type Stats struct{ Count int; Longest string }` stored in a `map[string]Stats` keyed by the first letter,
   and updates it without using pointers.
4. Bonus: start 4 goroutines that each call your grouping code on part of the input and merge into a shared map.
   Run it with `go run -race .` first without a mutex, observe the result, then fix it.

## Directory

`map_semantics/`
