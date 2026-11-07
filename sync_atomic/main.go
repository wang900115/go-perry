package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ---------- 1. Counter: atomic.Int64 ----------

type Stats struct {
	requests atomic.Int64 // zero value is ready to use, no init needed
	errors   atomic.Int64
}

// ---------- 2. CAS loop: lock-free "update only if larger" ----------

type MaxTracker struct {
	max atomic.Int64
}

func (m *MaxTracker) Observe(v int64) {
	for {
		old := m.max.Load()
		if v <= old {
			return // no update needed
		}
		if m.max.CompareAndSwap(old, v) {
			return // new value installed
		}
		// CAS failed: another goroutine changed it after our Load; reload and retry
	}
}

// ---------- 3. atomic.Pointer: config hot-swap (copy-on-write) ----------

type Config struct {
	Version  int
	Features []string
}

var current atomic.Pointer[Config]

func loadConfig() *Config { return current.Load() }

func reload(version int, features ...string) {
	// Build a brand-new Config and swap it in with a single Store; the GC reclaims the old one.
	// Never mutate the object returned by current.Load().
	current.Store(&Config{Version: version, Features: features})
}

// ---------- 4. Flag bits: Go 1.23 Or / And ----------

const (
	FlagReady uint32 = 1 << iota
	FlagDraining
	FlagClosed
)

func main() {
	fmt.Println("== 1. atomic.Int64 counter ==")
	var s Stats
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 1000 {
				s.requests.Add(1)
				if (i*1000+j)%100 == 0 {
					s.errors.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("requests=%d errors=%d\n", s.requests.Load(), s.errors.Load())

	fmt.Println("== 2. CAS loop max tracker ==")
	var mt MaxTracker
	for g := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range int64(10000) {
				mt.Observe((v*7919 + int64(g)*31) % 99991)
			}
		}()
	}
	wg.Wait()
	fmt.Println("max =", mt.max.Load())

	fmt.Println("== 3. atomic.Pointer config hot-swap ==")
	reload(1, "login")
	cfg := loadConfig() // readers get a snapshot of one complete version
	reload(2, "login", "dark-mode")
	fmt.Printf("old snapshot: v%d %v\n", cfg.Version, cfg.Features)
	fmt.Printf("current:      v%d %v\n", loadConfig().Version, loadConfig().Features)
	swapped := current.CompareAndSwap(cfg, &Config{Version: 99})
	fmt.Println("CAS with stale pointer succeeded?", swapped)

	fmt.Println("== 4. Flag bits Or / And (Go 1.23+) ==")
	var state atomic.Uint32
	old := state.Or(FlagReady | FlagDraining) // returns the OLD value
	fmt.Printf("Or  before=%03b after=%03b\n", old, state.Load())
	old = state.And(^FlagDraining) // clear Draining
	fmt.Printf("And before=%03b after=%03b\n", old, state.Load())

	fmt.Println("== 5. atomic.Value panics on inconsistent type ==")
	var v atomic.Value
	v.Store("hello")
	func() {
		defer func() { fmt.Println("recover:", recover()) }()
		v.Store(42)
	}()
	fmt.Println("value is still:", v.Load())
}
