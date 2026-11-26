# container/heap: Priority Queues in the Standard Library

> Go version used: `go1.26.3`. `container/heap` has been stable since Go 1.0; this lesson uses `any` (Go 1.18+).

## Why learn this

Whenever you need "give me the most urgent / smallest / earliest item next" while items keep arriving:

- Job schedulers and retry queues ordered by next-run time
- Timer wheels / deadline management (the Go runtime keeps each P's timers in its own 4-ary heap, `timerHeapN = 4` in `runtime/time.go`; it does not use `container/heap`)
- Dijkstra / A* shortest path, merging K sorted streams
- Top-K queries ("10 slowest requests"), streaming median (two heaps)
- LRU/LFU-style eviction by score

Sorting the whole slice after every insert is O(n log n) per insert. A binary heap gives O(log n) insert and O(log n) remove-min.

## Core concepts

### 1. The interface

```go
type Interface interface {
    sort.Interface        // Len() int; Less(i, j int) bool; Swap(i, j int)
    Push(x any)           // add x as element Len()
    Pop() any             // remove and return element Len() - 1
}
```

Docs: https://pkg.go.dev/container/heap

The package does **not** give you a ready-made heap type. You supply a slice-backed type, and the package
functions `heap.Init`, `heap.Push`, `heap.Pop`, `heap.Fix`, `heap.Remove` do the algorithm on it.

### 2. How it works

The heap is stored as an implicit binary tree in a slice: for index `i`, children are `2i+1` and `2i+2`, parent is `(i-1)/2`.
Invariant: `!h.Less(child, parent)` for every node, so `h[0]` is always the "smallest" according to `Less`.

Reading the source (https://cs.opensource.google/go/go/+/refs/tags/go1.26.3:src/container/heap/heap.go):

| Function | What it does | Cost |
| --- | --- | --- |
| `heap.Init(h)` | sift down from `n/2-1` to `0` (bottom-up heapify) | O(n) |
| `heap.Push(h, x)` | `h.Push(x)` (your append), then **up** from the last index | O(log n) |
| `heap.Pop(h)` | `Swap(0, n-1)`, **down** from 0 over the first n-1 items, then `h.Pop()` (your remove-last) | O(log n) |
| `heap.Remove(h, i)` | `Swap(i, n-1)`, down-or-up at `i`, then `h.Pop()` | O(log n) |
| `heap.Fix(h, i)` | element `i` changed: down, and if it did not move, up | O(log n) |

So **your** `Push` / `Pop` methods are only the "grow / shrink the slice" hooks. The ordering work is done by the package.
This is why they need pointer receivers (they change the slice header), while `Len/Less/Swap` can use value receivers.

The package works through an interface, so every `Less`/`Swap` is a dynamic call, and `Push(any)` boxes values.
For very hot paths people sometimes hand-write a typed heap; there is no generic `container/heap` in the standard library as of Go 1.26 (I am not aware of an accepted proposal that has landed).

### 3. Priority queue with index tracking

To change an element's priority or delete an arbitrary element, you need its current position. Store an `index` field and keep it
updated in `Swap`, `Push`, and `Pop` (exactly the pattern in the official `PriorityQueue` example:
https://pkg.go.dev/container/heap#example-package-PriorityQueue). Then:

```go
task.priority = 10
heap.Fix(pq, task.index)        // re-establish order
heap.Remove(pq, other.index)    // delete arbitrary element
```

### 4. Min vs max, ties, top-K

- `Less` with `<` = min-heap, `>` = max-heap. Nothing else changes.
- A heap is **not stable**: equal elements come out in no particular order. Add a monotonically increasing `seq` and compare it when priorities tie.
- Top-K largest of a stream: keep a **min**-heap of size K. If a new value is larger than `h[0]`, overwrite `h[0]` and call `heap.Fix(h, 0)`. O(n log K) time, O(K) memory.

## Complete example

`container_heap/main.go` (run with `go run .`):

```go
// Package main demonstrates container/heap: the heap.Interface contract,
// a priority queue with heap.Fix / heap.Remove, the "call h.Push directly"
// pitfall, stable ordering with a sequence number, and a bounded top-K heap.
package main

import (
	"container/heap"
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// 1. IntHeap: the smallest possible heap.Interface implementation.
// ---------------------------------------------------------------------------

// IntHeap is a min-heap of ints.
// sort.Interface (Len/Less/Swap) can use a value receiver,
// but Push/Pop change the slice length, so they need a pointer receiver.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // < means min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push is called BY heap.Push after nothing else: just append at the end.
func (h *IntHeap) Push(x any) { *h = append(*h, x.(int)) }

// Pop is called BY heap.Pop after it has swapped the root to the end:
// just remove and return the LAST element.
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func demoIntHeap() {
	fmt.Println("== 1. IntHeap: Init / Push / Pop ==")
	h := &IntHeap{5, 2, 8, 1, 9}
	heap.Init(h) // O(n) heapify
	fmt.Println("after Init  :", *h, "(root h[0] is the minimum)")

	heap.Push(h, 3) // O(log n): append + sift up
	fmt.Println("after Push 3:", *h)

	var out []int
	for h.Len() > 0 {
		out = append(out, heap.Pop(h).(int)) // O(log n): swap root to end, sift down, call h.Pop
	}
	fmt.Println("pop order   :", out)
}

// ---------------------------------------------------------------------------
// 2. Priority queue with index tracking, heap.Fix and heap.Remove.
// ---------------------------------------------------------------------------

type Task struct {
	name     string
	priority int // higher = more urgent
	index    int // position in the heap, maintained by Swap/Push/Pop
}

type PQ []*Task

func (pq PQ) Len() int { return len(pq) }

// Less with > turns it into a max-heap on priority.
func (pq PQ) Less(i, j int) bool { return pq[i].priority > pq[j].priority }

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PQ) Push(x any) {
	t := x.(*Task)
	t.index = len(*pq)
	*pq = append(*pq, t)
}

func (pq *PQ) Pop() any {
	old := *pq
	n := len(old)
	t := old[n-1]
	old[n-1] = nil // let the GC reclaim the task once the caller drops it
	t.index = -1   // mark as "not in the heap"
	*pq = old[:n-1]
	return t
}

func (pq PQ) String() string {
	parts := make([]string, len(pq))
	for i, t := range pq {
		parts[i] = fmt.Sprintf("%s:%d@%d", t.name, t.priority, t.index)
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func demoPriorityQueue() {
	fmt.Println("\n== 2. Priority queue: heap.Fix / heap.Remove ==")
	tasks := map[string]*Task{}
	pq := &PQ{}
	for _, t := range []struct {
		name string
		prio int
	}{{"email", 2}, {"deploy", 5}, {"lunch", 1}, {"review", 3}} {
		task := &Task{name: t.name, priority: t.prio}
		tasks[t.name] = task
		heap.Push(pq, task)
	}
	fmt.Println("heap        :", pq)

	// Change a priority in place, then tell the heap where it changed.
	lunch := tasks["lunch"]
	lunch.priority = 10
	heap.Fix(pq, lunch.index) // O(log n), cheaper than Remove + Push
	fmt.Println("lunch -> 10 :", pq)

	// Remove an arbitrary element by its index.
	removed := heap.Remove(pq, tasks["email"].index).(*Task)
	fmt.Printf("removed     : %s (index now %d)\n", removed.name, removed.index)

	for pq.Len() > 0 {
		t := heap.Pop(pq).(*Task)
		fmt.Printf("pop         : %s (priority %d)\n", t.name, t.priority)
	}
}

// ---------------------------------------------------------------------------
// 3. Pitfall: calling the methods directly instead of the heap functions.
// ---------------------------------------------------------------------------

func isMinHeap(h IntHeap) bool {
	for i := 1; i < len(h); i++ {
		if h[(i-1)/2] > h[i] { // parent of i is (i-1)/2
			return false
		}
	}
	return true
}

func demoDirectCallPitfall() {
	fmt.Println("\n== 3. Pitfall: h.Push / h.Pop vs heap.Push / heap.Pop ==")
	h := &IntHeap{2, 4, 6}
	heap.Init(h)

	h.Push(1) // WRONG: only appends, no sift up
	fmt.Println("h.Push(1)       :", *h, "valid heap?", isMinHeap(*h))
	fmt.Println("heap.Pop(h)     :", heap.Pop(h), "<- should be 1, the broken invariant returns 2")

	h2 := &IntHeap{2, 4, 6}
	heap.Init(h2)
	heap.Push(h2, 1) // RIGHT
	fmt.Println("heap.Push(h2, 1):", *h2, "valid heap?", isMinHeap(*h2))
	fmt.Println("heap.Pop(h2)    :", heap.Pop(h2))
}

// ---------------------------------------------------------------------------
// 4. A heap is not stable: break ties with a sequence number.
// ---------------------------------------------------------------------------

type Job struct {
	name string
	prio int
	seq  int // insertion order
}

type JobHeap struct {
	items  []Job
	stable bool
}

func (h *JobHeap) Len() int { return len(h.items) }
func (h *JobHeap) Less(i, j int) bool {
	a, b := h.items[i], h.items[j]
	if a.prio != b.prio {
		return a.prio > b.prio
	}
	if h.stable {
		return a.seq < b.seq // FIFO among equal priorities
	}
	return false
}
func (h *JobHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *JobHeap) Push(x any)    { h.items = append(h.items, x.(Job)) }
func (h *JobHeap) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}

func drainJobs(stable bool) string {
	h := &JobHeap{stable: stable}
	for i, name := range []string{"A", "B", "C", "D", "E"} {
		heap.Push(h, Job{name: name, prio: 1, seq: i}) // all the same priority
	}
	var out []string
	for h.Len() > 0 {
		out = append(out, heap.Pop(h).(Job).name)
	}
	return strings.Join(out, " ")
}

func demoStability() {
	fmt.Println("\n== 4. Equal priorities: heap order is not insertion order ==")
	fmt.Println("without seq:", drainJobs(false))
	fmt.Println("with seq   :", drainJobs(true))
}

// ---------------------------------------------------------------------------
// 5. Top-K largest with a bounded MIN-heap of size K: O(n log k).
// ---------------------------------------------------------------------------

func topK(nums []int, k int) []int {
	h := &IntHeap{}
	for _, n := range nums {
		if h.Len() < k {
			heap.Push(h, n)
		} else if n > (*h)[0] { // bigger than the smallest of the current top-K
			(*h)[0] = n
			heap.Fix(h, 0) // replace root and sift down in one step
		}
	}
	// Popping a min-heap yields ascending order; reverse for descending.
	res := make([]int, h.Len())
	for i := len(res) - 1; i >= 0; i-- {
		res[i] = heap.Pop(h).(int)
	}
	return res
}

func demoTopK() {
	fmt.Println("\n== 5. Top-3 largest with a size-3 min-heap ==")
	nums := []int{7, 1, 9, 3, 9, 4, 12, 0, 5}
	fmt.Println("input :", nums)
	fmt.Println("top 3 :", topK(nums, 3))
}

func main() {
	demoIntHeap()
	demoPriorityQueue()
	demoDirectCallPitfall()
	demoStability()
	demoTopK()
}
```

Expected output (verified with `go1.26.3`):

```text
== 1. IntHeap: Init / Push / Pop ==
after Init  : [1 2 8 5 9] (root h[0] is the minimum)
after Push 3: [1 2 3 5 9 8]
pop order   : [1 2 3 5 8 9]

== 2. Priority queue: heap.Fix / heap.Remove ==
heap        : [deploy:5@0 review:3@1 lunch:1@2 email:2@3]
lunch -> 10 : [lunch:10@0 review:3@1 deploy:5@2 email:2@3]
removed     : email (index now -1)
pop         : lunch (priority 10)
pop         : deploy (priority 5)
pop         : review (priority 3)

== 3. Pitfall: h.Push / h.Pop vs heap.Push / heap.Pop ==
h.Push(1)       : [2 4 6 1] valid heap? false
heap.Pop(h)     : 2 <- should be 1, the broken invariant returns 2
heap.Push(h2, 1): [1 2 6 4] valid heap? true
heap.Pop(h2)    : 1

== 4. Equal priorities: heap order is not insertion order ==
without seq: A E D C B
with seq   : A B C D E

== 5. Top-3 largest with a size-3 min-heap ==
input : [7 1 9 3 9 4 12 0 5]
top 3 : [12 9 9]
```

Note: the exact internal layout printed after `Init`/`Push` and the "without seq" order are deterministic for this
implementation, but they are implementation details. Only the pop order (by priority) is guaranteed.

## Common pitfalls

### Pitfall 1: calling `h.Push` / `h.Pop` instead of `heap.Push` / `heap.Pop`

Your methods only append / remove the last element. `h.Push(1)` leaves `[2 4 6 1]`, which is not a heap, and the next
`heap.Pop` returns `2` instead of `1`. Always go through the package functions (section 3 of the example).

### Pitfall 2: changing priority without `heap.Fix` (or forgetting to maintain `index`)

Mutating `task.priority` silently breaks the invariant; nothing panics, you just get wrong order later.
If `Swap` does not update `index`, `heap.Fix(pq, task.index)` fixes the wrong slot. After `Pop`, set `index = -1` so stale use is visible,
and set the vacated slot to `nil` so the backing array does not keep the object alive.

### Pitfall 3: expecting FIFO among equal priorities / expecting sorted slice contents

The heap slice is only partially ordered: `h[0]` is the min, `h[1]` is **not** necessarily the second smallest.
Do not range over the slice expecting order; pop instead (or sort a copy). And equal priorities do not come out in insertion order
(`A E D C B` in the example) unless you add a tie-breaker.

## Exercise (no answer given)

Write a function `mergeK(lists [][]int) []int` that merges K **already sorted** int slices into one sorted slice using `container/heap`.

Requirements:
1. The heap must hold at most K items at any time (one "cursor" per list), not all elements.
2. Use `heap.Fix(h, 0)` instead of `Pop` + `Push` when advancing a cursor whose list still has items.
3. Handle empty input lists.
4. Test with `[][]int{{1, 4, 9}, {}, {2, 3, 10, 11}, {5}}` and state the time complexity in terms of N (total elements) and K.

Bonus: make it stable — when two lists have the same value, the one with the smaller list index must come first.

## Suggested directory

`container_heap/`
