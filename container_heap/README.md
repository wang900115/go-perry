# container_heap

Priority queues with the standard library `container/heap` package.

- `heap.Interface` = `sort.Interface` (`Len`/`Less`/`Swap`) + `Push(x any)` + `Pop() any`.
- `Less` decides the order: `<` gives a min-heap, `>` gives a max-heap.
- Always call the package functions `heap.Init/Push/Pop/Fix/Remove`; your `Push`/`Pop` methods only append / remove the last element.
- Track each element's `index` in `Swap`/`Push`/`Pop` so you can use `heap.Fix` and `heap.Remove`.
- A heap is not stable: add a sequence number to break ties if you need FIFO among equal priorities.
- Top-K largest: keep a size-K min-heap and replace the root (`h[0] = x; heap.Fix(h, 0)`).

```bash
go run .
```

See [LESSON.md](LESSON.md) for the full lesson.
