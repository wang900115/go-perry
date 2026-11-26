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
