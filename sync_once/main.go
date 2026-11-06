package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// ---------- 1. sync.Once：多個 goroutine 同時呼叫，init 只跑一次 ----------

type Config struct {
	DSN     string
	MaxConn int
}

var (
	cfgOnce   sync.Once
	cfg       *Config
	loadCount atomic.Int32
)

func GetConfig() *Config {
	cfgOnce.Do(func() {
		loadCount.Add(1)
		cfg = &Config{DSN: "file:app.db", MaxConn: 10}
	})
	return cfg
}

func demoOnce() {
	fmt.Println("== 1. sync.Once ==")
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = GetConfig()
		}()
	}
	wg.Wait()
	fmt.Printf("100 個 goroutine 呼叫 GetConfig，實際載入次數 = %d\n", loadCount.Load())
	fmt.Printf("config = %+v\n", *GetConfig())
}

// ---------- 2. OnceValue / OnceValues：把「只算一次 + 回傳值」包成一個函式 ----------

var expensiveTable = sync.OnceValue(func() []int {
	fmt.Println("  (計算平方表中…只會印一次)")
	t := make([]int, 10)
	for i := range t {
		t[i] = i * i
	}
	return t
})

var loadSecret = sync.OnceValues(func() (string, error) {
	fmt.Println("  (讀取 secret 中…只會印一次)")
	return "", errors.New("secret file not found")
})

func demoOnceValue() {
	fmt.Println("== 2. OnceValue / OnceValues ==")
	fmt.Println("table[3] =", expensiveTable()[3])
	fmt.Println("table[9] =", expensiveTable()[9])

	for i := 1; i <= 2; i++ {
		s, err := loadSecret()
		fmt.Printf("第 %d 次 loadSecret: %q, err=%v\n", i, s, err)
	}
}

// ---------- 3. panic 行為差異：Once.Do vs OnceFunc ----------

func callSafely(name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  %s -> panic: %v\n", name, r)
			return
		}
		fmt.Printf("  %s -> 正常返回（f 沒有再被執行）\n", name)
	}()
	f()
}

func demoPanic() {
	fmt.Println("== 3. panic 時的行為 ==")

	var once sync.Once
	boom := func() { panic("init failed") }
	callSafely("Once.Do 第 1 次", func() { once.Do(boom) })
	callSafely("Once.Do 第 2 次", func() { once.Do(boom) })

	of := sync.OnceFunc(boom)
	callSafely("OnceFunc 第 1 次", of)
	callSafely("OnceFunc 第 2 次", of)
}

func main() {
	demoOnce()
	demoOnceValue()
	demoPanic()
}
