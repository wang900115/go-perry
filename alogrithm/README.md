# alogrithm

Hand-written data structures, search routines, and LeetCode-style "design" problems (LRU/LFU caches, stream problems). It is a work in progress: several files are stubs, and some implementations have bugs, listed under Notes.

## What the code does

There is no `main` to run. Everything is a library package.

### `design/cache/`
- `lru/lru.go` - `LRUCache`: a hash map plus a doubly linked list with sentinel head and tail nodes. `Get` and `Put` move the node to the front, and a full cache evicts `tail.prev`. It contains a `func main` demo, but the package is `lru`, so that `main` never runs.
- `lfu/lfu.go` - `LFUCache`: `nodes` map, `freqMap` (frequency to doubly linked list), and `minFreq`. On access a node moves to the list for `freq+1`. Eviction removes the tail (least recent) of the `minFreq` list.

### `design/stream/`
- `char/char.go` - `SteamCheck` (stream of characters): builds a trie of the **reversed** words. `Query(letter)` appends to a buffer capped at the longest word length, then walks the trie backwards from the newest letter. It returns true if some word is a suffix of the stream.
- `disjoint/disjoint.go` - `SummaryRanges`: stores intervals as `start to end` in a map. `AddNum` merges with an adjacent interval on either side. `GetIntervals` sorts the starts.
- `kth/kth.go` - `KthLargest`: a size-k min-heap. The root is the k-th largest value, and `Add` replaces the root when the new value is at least that large.

### `search/` (package `search`)
- `linear.go` - `linear(arr, target)`: unexported O(n) scan.
- `binary.go` - `Binary(arr, target)`: binary search with an overflow-safe `mid`.
- `bfs.go` - `BFS(graph, start)`: queue-based traversal over an adjacency map.
- `dfs.go` - `DFS` (recursive) and `DFS_` (explicit stack, neighbors pushed in reverse so they pop in order).

The BFS/DFS functions only mark nodes as visited. They return and print nothing.

### `structure/`
- `binarytree.go` (package `binarytree`) - binary search tree: `Insert` (rejects duplicates), `Search`, `Delete`, `Min`, `Max`. `InOrder`/`PreOrder`/`PostOrder` are empty.
- `graph.go` (package `graph`) - adjacency-list `Graph` with `AddEdge` (adds both directions unless `directed`), `Neighbors`, `Exist`.
- `heap.go` (package `heap`) - int min-heap: `Insert`, `Remove`, `Peek`, `Size`, `HeapifyUp`/`HeapifyDown`.
- `linkedlist.go` (package `linkedlist`) - singly linked list: `Append`, `Prepend`, `Delete`, `Find`.
- `queue.go` / `stack.go` - slice-backed queue and stack.

### `sort/`
`bubble.go`, `heap.go`, `insertion.go`, `merge.go`, `quick.go`, `radix.go`, `selection.go` are **empty stubs** (only `package sort`).

## Concepts you can learn

- **Hash map + doubly linked list (LRU)**: the map finds a node in O(1), and the list reorders or evicts it in O(1). Sentinel head and tail nodes remove nil checks at the edges.
- **Frequency buckets (LFU)**: tracking `minFreq` makes eviction O(1) instead of scanning every frequency. Ties within a frequency fall back to LRU order.
- **Reversed trie for suffix matching**: matching words that end at the newest character only needs a backward walk bounded by the longest word, O(maxLen) per query.
- **Size-k min-heap for top-k**: O(log k) per insert with O(k) memory, instead of sorting the whole stream.
- **BFS vs DFS**: a FIFO queue explores level by level, while a stack or recursion goes deep first. Recursive DFS can overflow the stack on deep graphs, and the explicit stack avoids that.
- **Binary search**: O(log n) on sorted input. `left + (right-left)/2` avoids integer overflow.
- **BST deletion cases** (0, 1, or 2 children, using the in-order successor), and **array-backed binary heaps** (parent `(i-1)/2`, children `2i+1` / `2i+2`).

## How to run

Part of the root module (no `go.mod` here), and there is no `main` anywhere, so there is nothing to `go run`. There are no tests either. To type-check, run from the repo root:

```bash
go vet ./alogrithm/design/... ./alogrithm/search ./alogrithm/sort
go vet alogrithm/structure/heap.go        # structure/: one file at a time
```

To use a package, import it from another root-module program, for example `import "arena/alogrithm/design/cache/lfu"`.

## Notes

Problems visible in the code:

- **`structure/` does not build as a directory**: each file declares a different package (`binarytree`, `graph`, `heap`, ...), so `go vet ./alogrithm/...` fails with "found packages binarytree and graph". Each file only type-checks on its own.
- **`search/binary.go` comparison is inverted**: it moves `left` up when `arr[mid] > target`, so it only works on arrays sorted in **descending** order.
- **`lru.go` `Put` panics on a new key**: it calls `moveToHead` on a fresh node whose `prev`/`next` are nil (nil pointer dereference). The in-file `main` demo would hit this.
- **`binarytree.go` `Delete`**: `parent` is never updated while searching, so deleting a non-root node with 0 or 1 child overwrites `b.root`. In the two-children case, the successor unlink sets the wrong side (`.right` when the successor is a `.left` child, and the reverse).
- **`linkedlist.go` `Prepend`**: sets `newNode.next = l.head.next`, which drops the old head (and panics on an empty list). It should be `l.head`.
- `graph.Graph` has no constructor, and a zero-value `Graph` has a nil map, so `AddEdge` panics on it.
- `char.go` assumes lowercase `a`–`z` (`b - 'a'` indexes a `[26]` array), and the type is spelled `SteamCheck`.
- `disjoint.AddNum` scans all intervals, so each insert is O(n), not O(log n).
- The directory name `alogrithm` is misspelled but real. Do not rename it.
