# 2026-11-18 · Slice internals: header, aliasing, growth

## Why learn this
- Slices are the most used data structure in Go. Every `[]byte` buffer, every list returned from a DB query, every `bufio.Scanner.Bytes()` is a slice.
- Many real bugs come from not knowing that slices share memory: two `append` calls overwriting each other, a function whose `append` "disappears", a 16-byte header keeping a 1 MiB buffer alive, or `slices.Delete` changing a slice you still hold.
- It explains earlier lessons: why `Scanner.Bytes()` gets overwritten (`bufio_scanner`), and why `io.Reader.Read(p)` fills a caller-owned slice (`io_reader_writer`).

## Core concepts
1. **A slice is a 3-word header**: pointer to an element of a backing array, `len`, and `cap` (https://go.dev/ref/spec#Slice_types, https://go.dev/blog/slices-intro). Copying a slice copies the header, not the elements.
2. **Slicing shares memory**: `a[lo:hi]` points into the same array; `cap` is counted from `lo` to the end of the array. You may re-slice up to `cap`, beyond `len`.
3. **`append` writes in place when `len < cap`**. Only when `cap` is too small does it allocate a new array and copy. So the result may or may not alias the input; always use the return value: `s = append(s, v)`.
4. **Full slice expression** `a[lo:hi:max]` sets `cap = max-lo` (https://go.dev/ref/spec#Slice_expressions). `slices.Clip(s)` is `s[:len(s):len(s)]`. A later `append` must then copy.
5. **Passing to a function copies the header**: element writes are visible to the caller; a change of `len` (append) is not. Return the new slice, like `append` does.
6. **`slices.Delete`, `DeleteFunc`, `Compact`, `CompactFunc`, `Replace` zero the elements between the new and old length since Go 1.22** (https://go.dev/doc/go1.22#slices, https://pkg.go.dev/slices#Delete). This avoids leaking pointers but changes what old headers see.
7. **Sub-slices keep the whole array alive**: the GC sees one pointer into the array, so the whole array stays reachable. Copy the small part with `slices.Clone` (or `bytes.Clone`).
8. **nil vs empty**: both have `len 0` and work with `append`, `range`, `len`; only `s == nil` differs (and, for example, `encoding/json` encodes them as `null` vs `[]`).

### How it works (runtime)
Based on go1.26 `src/runtime/slice.go` (`growslice`, `nextslicecap`); these are implementation details, not part of the spec, and may change.
- When `append` needs more room, the compiler calls `runtime.growslice`. `nextslicecap` doubles the capacity while it is below 256 elements; above that it grows by about `newcap += (newcap + 3*256) / 4`, a smooth move from 2x to about 1.25x.
- The byte size is then rounded up to a malloc size class (`roundupsize`), so the final `cap` can be larger than the formula. Example from Part 3: 512 → 832 by the formula; 832×8 = 6656 bytes is rounded to the 6784-byte class, so `cap = 848`.
- For element types that contain no pointers, only the copied part is cleared/copied; for pointer types the new memory is zeroed so the GC never sees garbage.
- The old array is not freed by `append`; it is collected later if nothing else references it.

## Full example
See `main.go`. Verified with `go run .` (2 runs, identical output) and `go vet .` on go1.26 linux/amd64; `GOOS=windows go vet .` also passes. Expected output is in [README.md](README.md).

Parts:
1. Sub-slices of an array, writing through one and seeing it in another, re-slicing up to `cap`.
2. Two `append` calls to the same base overwrite each other; the `s[:3:3]` / `slices.Clip` fix.
3. Capacity sequence while appending one element at a time; growth detaches the new slice from the old one; preallocation.
4. Element writes vs `append` inside a function.
5. `slices.Delete` zeroing the tail, seen through the old header; no leaked pointers.
6. A 16-byte sub-slice pinning 1 MiB, `slices.Clone`, nil vs empty slices.

## Common pitfalls
1. **Appending to a shared slice.** `a := append(base, 1); b := append(base, 2)` - if `base` has spare capacity, `a` and `b` share the same last element. Typical in recursion/backtracking (`path = append(path, x)` then storing `path`). Store `slices.Clone(path)` or use `base[:len(base):len(base)]`.
2. **Using a slice after `slices.Delete` (or forgetting the return value).** The old header keeps the old `len` and now sees zeros (or `nil` pointers). Always write `s = slices.Delete(s, i, j)` and do not keep other headers to the same array.
3. **Keeping a small piece of a big buffer.** Returning `buf[:16]` from a 1 MiB read buffer, or storing `line := data[a:b]` in a long-lived map, keeps the whole buffer in memory. Clone what you keep.

## Exercise (no answer given)
Write a `Stack[T any]` type backed by a slice with `Push`, `Pop` and `Peek`:
- Make `Pop` not leak memory when `T` is a pointer type (hint: what should happen to the popped slot?). Prove it with a test that uses `runtime.AddCleanup` or `weak.Pointer` (see the `weak_pointer` topic).
- Add a method `Snapshot() []T`. Write a test showing that a naive `return s.items` lets the caller corrupt the stack after the next `Push`, then fix it.
- Write a function `subsets(nums []int) [][]int` with backtracking. First write it the naive way (`res = append(res, path)`), find the input size where results start to be wrong, and explain using `len`/`cap`. Then fix it.
- Bonus: print the `cap` sequence for `[]byte` and for `[]struct{ a, b, c int64 }` and explain why it differs from `[]int` (size classes).

## Suggested directory
`slice_internals`
