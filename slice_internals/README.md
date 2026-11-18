# slice_internals

How Go slices really work: the slice header (pointer, `len`, `cap`), sub-slices
sharing one backing array, the `append` aliasing bug and the full slice
expression `s[lo:hi:max]` / `slices.Clip` fix, growth rules in `runtime.growslice`,
why a function must return the slice it appended to, Go 1.22 `slices.Delete`
zeroing the tail, and memory retention by small sub-slices (`slices.Clone`).

```bash
go run .
```

Standard library only. Verified on go1.26 linux/amd64 (2 runs, identical
output); `go vet .` and `GOOS=windows go vet .` pass. The exact capacities in
Part 3 (and `cap=6` for `grown`) come from runtime growth rules and allocator
size classes, so they can differ between Go versions.
See [LESSON.md](LESSON.md) for the full lesson.

## Expected output

```text
== Part 1: slice header and shared backing array ==
a      len=3 cap=5 [1 2 3]
b      len=1 cap=4 [2]
after b[0] = 99:
a      len=3 cap=5 [1 99 3]
arr   [0 1 99 3 4 5]
b[:4]  len=4 cap=4 [99 3 4 5]

== Part 2: append aliasing ==
x      len=4 cap=10 [1 2 3 200]
y      len=4 cap=10 [1 2 3 200]
x[3] changed by appending to y: true
x2     len=4 cap=6 [1 2 3 100]
y2     len=4 cap=6 [1 2 3 200]
cap(slices.Clip(base)) = 3

== Part 3: growth (cap changes while appending 1 by 1) ==
caps: [1 2 4 8 16 32 64 128 256 512 848 1280 1792 2560]
old    len=3 cap=3 [1 2 3]
grown  len=4 cap=6 [42 2 3 4]
preallocated: len 2000 cap 2000

== Part 4: passing slices to functions ==
set    len=2 cap=4 [-1 0]
app    len=2 cap=4 [-1 0]
hidden element still in the array: [-1 0 7]
ret    len=3 cap=4 [-1 0 7]

== Part 5: slices.Delete / Compact zero the tail (Go 1.22+) ==
d      len=3 cap=5 [1 4 5]
s      len=5 cap=5 [1 4 5 0 0]
len(users2) = 2 users[2] == nil: true (no leaked pointer)

== Part 6: memory retention by sub-slices ==
head: len=16 cap=1048576 (pins the 1 MiB array)
clone: len=16 cap=16 (new small array, big can be collected)
nil slice: true 0 | empty slice: false 0
append to nil: [1]
```
