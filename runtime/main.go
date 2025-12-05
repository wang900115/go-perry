package main

import (
	"fmt"
	"runtime"
	"time"
)

func riskFunc() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			// 取得目前 goroutine 的 stack trace
			n := runtime.Stack(buf, false)
			println("Recovered from panic:", r)
			println("Stack trace:\n", string(buf[:n]))
		}
	}()

	var ptr *int

	fmt.Println(*ptr)
}

// 讓出 目前 goroutine 的 CPU 使用權，讓其他 goroutine 有機會執行
func goschedExample() {
	fmt.Println("Starting program...")
	go riskFunc()

	runtime.Gosched()
	fmt.Println("Program finished.")
}

// 設定最大可使用的 CPU 核心數量
func goMaxProcesExample() {
	prev := runtime.GOMAXPROCS(1)
	fmt.Println("Previous GOMAXPROCS:", prev)

	for i := 0; i < 5; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d is running\n", id)
			time.Sleep(1 * time.Second)
		}(i)
	}

	runtime.Gosched()
}

// 結束目前 goroutine 的執行
func exitExample() {
	go func() {
		defer fmt.Println("defer executed before Goexit") // 這會在Goexit前執行

		fmt.Println("goroutine exiting")
		runtime.Goexit()

		fmt.Println("This will not be printed")
	}()

	runtime.Gosched()
}

// 查看目前有多少個 goroutine 在執行
func seeGoroutine() {
	for i := 0; i < 5; i++ {
		go func() {
			time.Sleep(1 * time.Second)
		}()
	}

	fmt.Println("Current goroutines:", runtime.NumGoroutine())
}

type Resource struct {
	Name string
}

func fininalizerExample() {
	r := &Resource{Name: "MyResource"}

	// 設定終結器
	runtime.SetFinalizer(r, func(res *Resource) {
		fmt.Println("Finalizer called for:", res.Name)
	})

	r = nil // 解除引用，讓 GC 可以回收

	// 強制 GC
	runtime.GC()

	// 等待一點時間讓終結器執行
	runtime.Gosched()
}
