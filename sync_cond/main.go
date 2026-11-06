package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// BoundedQueue 是一個用 sync.Cond 實作的有界佇列（容量固定的 FIFO）。
// 一把 Mutex 保護狀態，兩個 Cond 共用這把鎖：
//
//	notEmpty：佇列空時，消費者在這裡等
//	notFull ：佇列滿時，生產者在這裡等
type BoundedQueue struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	items    []int
	cap      int
	closed   bool
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{cap: capacity}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

func (q *BoundedQueue) Put(v int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// 一定要用 for，不能用 if：被喚醒後條件不一定成立
	for len(q.items) == q.cap {
		q.notFull.Wait() // 原子地：解鎖 → 睡覺 → 被喚醒後重新上鎖
	}
	q.items = append(q.items, v)
	q.notEmpty.Signal() // 叫醒「一個」等資料的消費者
}

// Get 回傳 (值, true)；佇列已關閉且清空時回傳 (0, false)。
func (q *BoundedQueue) Get() (int, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.items) == 0 && !q.closed {
		q.notEmpty.Wait()
	}
	if len(q.items) == 0 { // closed 且沒資料了
		return 0, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	q.notFull.Signal()
	return v, true
}

func (q *BoundedQueue) Close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
	q.notEmpty.Broadcast() // 叫醒「所有」消費者，讓他們看到 closed 然後離開
}

func demoBroadcast() {
	fmt.Println("== Part 1: Broadcast 一次叫醒全部 ==")
	var (
		mu    sync.Mutex
		cond  = sync.NewCond(&mu)
		ready bool
		wg    sync.WaitGroup
		order []string
	)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			for !ready {
				cond.Wait()
			}
			order = append(order, fmt.Sprintf("runner-%d", id))
			mu.Unlock()
		}(i)
	}
	time.Sleep(50 * time.Millisecond) // 只是讓 runner 有機會先進入 Wait，示範用
	mu.Lock()
	ready = true
	mu.Unlock()
	cond.Broadcast()
	wg.Wait()
	sort.Strings(order) // 起跑順序不保證，排序後輸出才穩定
	fmt.Println("all started:", order)
}

func demoQueue() {
	fmt.Println("== Part 2: 用兩個 Cond 做有界佇列 ==")
	q := NewBoundedQueue(2)

	var wg sync.WaitGroup
	results := make([]int, 3) // 每個 consumer 各自累加
	counts := make([]int, 3)
	for c := 0; c < 3; c++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				v, ok := q.Get()
				if !ok {
					return
				}
				results[id] += v
				counts[id]++
			}
		}(c)
	}

	for i := 1; i <= 10; i++ {
		q.Put(i) // 容量只有 2，生產者常會在 notFull 上等待
	}
	q.Close()
	wg.Wait()

	total, n := 0, 0
	for i := range results {
		total += results[i]
		n += counts[i]
	}
	fmt.Printf("consumed %d items, sum = %d\n", n, total)
}

func main() {
	demoBroadcast()
	demoQueue()
}
