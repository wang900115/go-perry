package main

import (
	"fmt"
	"runtime"
	"weak"
)

// Blob 模擬一個「很大、可以重新產生」的資源，例如解碼後的圖片。
type Blob struct {
	Name string
	Data []byte
}

// WeakCache：只持有弱指標，不會阻止 GC 回收 value。
type WeakCache struct {
	m map[string]weak.Pointer[Blob]
}

func NewWeakCache() *WeakCache {
	return &WeakCache{m: make(map[string]weak.Pointer[Blob])}
}

func (c *WeakCache) Put(key string, b *Blob) {
	c.m[key] = weak.Make(b)
	// 物件被回收時，順便把 map 裡的過期 entry 清掉，避免 map 無限長大。
	runtime.AddCleanup(b, func(k string) {
		// 注意：cleanup 在另一個 goroutine 執行；正式程式要加鎖。
		// 這裡為了示範簡單，只印出訊息。
		fmt.Println("  [cleanup] blob 已被回收:", k)
	}, key)
}

func (c *WeakCache) Get(key string) (*Blob, bool) {
	wp, ok := c.m[key]
	if !ok {
		return nil, false
	}
	b := wp.Value() // 若已被回收，回傳 nil
	return b, b != nil
}

func main() {
	cache := NewWeakCache()

	// 1. 放入兩個 blob，其中 a 由 main 持有強參照
	a := &Blob{Name: "a", Data: make([]byte, 1<<20)}
	cache.Put("a", a)
	cache.Put("b", &Blob{Name: "b", Data: make([]byte, 1<<20)}) // 沒有人強參照 b

	// 2. 同一個物件做出的 weak.Pointer 可以用 == 比較
	fmt.Println("weak.Make(a) == weak.Make(a):", weak.Make(a) == weak.Make(a))

	// 3. 觸發 GC
	runtime.GC()
	runtime.GC() // 第二次讓 cleanup goroutine 有機會跑完（示範用）

	if b, ok := cache.Get("a"); ok {
		fmt.Println("Get(a): 命中,", b.Name, len(b.Data), "bytes")
	}
	if _, ok := cache.Get("b"); !ok {
		fmt.Println("Get(b): 已被 GC 回收 → 需要重新產生")
	}

	// 4. 放掉 a 的最後一個強參照
	runtime.KeepAlive(a) // 確保 a 活到這一行
	a = nil
	runtime.GC()
	runtime.GC()

	if _, ok := cache.Get("a"); !ok {
		fmt.Println("Get(a): 強參照消失後，也被回收了")
	}
}
