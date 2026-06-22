package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestReadOrDone(t *testing.T) {
	t.Run("ReadOrDone", func(t *testing.T) {
		ctx := context.Background()
		in := make(chan int)
		go func() {
			in <- 1
			in <- 2
			in <- 3
			close(in)
		}()

		out := ReadOrDone(ctx, in)
		results := []int{}
		for v := range out {
			results = append(results, v)
		}

		if len(results) != 3 || results[0] != 1 || results[1] != 2 || results[2] != 3 {
			t.Errorf("expect [1 2 3], got %v", results)
		}
	})

	t.Run("context canceled stops reading", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		in := make(chan int)
		go func() {
			in <- 1
			in <- 2
			cancel() // cancel context
			in <- 3
			close(in)
		}()

		out := ReadOrDone(ctx, in)
		results := []int{}
		for v := range out {
			results = append(results, v)
		}

		if len(results) < 3 {
			t.Logf("context canceled stops reading, read %d values", len(results))
		}
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("FromSlice creates channel from slice", func(t *testing.T) {
		ctx := context.Background()
		slice := []int{1, 2, 3, 4, 5}
		out := FromSlice(ctx, slice)

		results := []int{}
		for v := range out {
			results = append(results, v)
		}

		if len(results) != len(slice) {
			t.Errorf("expect %d values, got %d", len(slice), len(results))
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		ctx := context.Background()
		out := FromSlice(ctx, []string{})

		count := 0
		for range out {
			count++
		}

		if count != 0 {
			t.Errorf("expect 0 values, got %d", count)
		}
	})
}

// TestMerge
func TestMerge(t *testing.T) {
	t.Run("merge multiple channels", func(t *testing.T) {
		ctx := context.Background()

		ch1 := FromSlice(ctx, []int{1, 2})
		ch2 := FromSlice(ctx, []int{3, 4})
		ch3 := FromSlice(ctx, []int{5, 6})

		merged := Merge(ctx, ch1, ch2, ch3)

		results := []int{}
		for v := range merged {
			results = append(results, v)
		}

		if len(results) != 6 {
			t.Errorf("expect 6 values, got %d", len(results))
		}
	})

	t.Run("context canceled stops merging", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		ch1 := FromSlice(ctx, []int{1, 2, 3, 4, 5})
		ch2 := FromSlice(ctx, []int{6, 7, 8, 9, 10})

		merged := Merge(ctx, ch1, ch2)

		count := 0
		for range merged {
			count++
			if count == 3 {
				cancel()
			}
		}

		t.Logf("context canceled stops merging, read %d values", count)
	})
}

func TestTee(t *testing.T) {
	t.Run("split into two channels", func(t *testing.T) {
		ctx := context.Background()
		in := FromSlice(ctx, []int{1, 2, 3})

		out1, out2 := Tee(ctx, in)

		results1 := []int{}
		for v := range out1 {
			results1 = append(results1, v)
		}

		results2 := []int{}
		for v := range out2 {
			results2 = append(results2, v)
		}

		if len(results1) != 3 || len(results2) != 3 {
			t.Errorf("expect both channels to have 3 values, got %d and %d", len(results1), len(results2))
		}

		for i := range results1 {
			if results1[i] != results2[i] {
				t.Errorf("values in both channels should be the same, got %d != %d", results1[i], results2[i])
			}
		}
	})

	t.Run("different consumption speeds for two channels", func(t *testing.T) {
		ctx := context.Background()
		in := FromSlice(ctx, []int{1, 2, 3, 4, 5})

		out1, out2 := Tee(ctx, in)

		// fast consumption of out1
		results1 := []int{}
		for v := range out1 {
			results1 = append(results1, v)
		}

		// slow consumption of out2
		results2 := []int{}
		for v := range out2 {
			results2 = append(results2, v)
		}

		if len(results1) != len(results2) {
			t.Errorf("both channels should have the same number of values, got %d and %d", len(results1), len(results2))
		}
	})
}

// TestStage
func TestStage(t *testing.T) {
	t.Run("process values and return results", func(t *testing.T) {
		ctx := context.Background()
		in := FromSlice(ctx, []int{1, 2, 3, 4, 5})

		// multiply each value by 2
		out, errCh := Stage(ctx, 2, in, func(ctx context.Context, v int) (int, error) {
			return v * 2, nil
		})

		results := []int{}
		errs := []error{}

		done := make(chan bool)
		go func() {
			for e := range errCh {
				errs = append(errs, e)
			}
			done <- true
		}()

		for v := range out {
			results = append(results, v)
		}
		<-done

		if len(results) != 5 {
			t.Errorf("expect 5 results, got %d", len(results))
		}
		if len(errs) != 0 {
			t.Errorf("expect 0 errors, got %d", len(errs))
		}
	})

	t.Run("error handling", func(t *testing.T) {
		ctx := context.Background()
		in := FromSlice(ctx, []int{1, 2, 3, 4, 5})

		// return error when value is 3
		out, errCh := Stage(ctx, 2, in, func(ctx context.Context, v int) (int, error) {
			if v == 3 {
				return 0, errors.New("processing failed for value 3")
			}
			return v * 2, nil
		})

		results := []int{}
		errs := []error{}

		done := make(chan bool)
		go func() {
			for e := range errCh {
				errs = append(errs, e)
			}
			done <- true
		}()

		for v := range out {
			results = append(results, v)
		}
		<-done

		if len(errs) == 0 {
			t.Errorf("expect errors, got 0")
		}
	})

	t.Run("concurrency control - max 2 workers", func(t *testing.T) {
		ctx := context.Background()
		in := FromSlice(ctx, []int{1, 2, 3, 4, 5})

		concurrentCount := 0
		maxConcurrent := 0
		mu := &sync.Mutex{}

		out, errCh := Stage(ctx, 2, in, func(ctx context.Context, v int) (int, error) {
			mu.Lock()
			concurrentCount++
			if concurrentCount > maxConcurrent {
				maxConcurrent = concurrentCount
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond) // 模拟处理时间

			mu.Lock()
			concurrentCount--
			mu.Unlock()

			return v, nil
		})

		results := []int{}
		errs := []error{}

		done := make(chan bool)
		go func() {
			for e := range errCh {
				errs = append(errs, e)
			}
			done <- true
		}()

		for v := range out {
			results = append(results, v)
		}
		<-done

		t.Logf("max concurrent workers: %d (expected <= 2)", maxConcurrent)
		if maxConcurrent > 2 {
			t.Errorf("concurrent workers exceeded limit: %d > 2", maxConcurrent)
		}
	})

	t.Run("stop processing when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		in := FromSlice(ctx, []int{1, 2, 3, 4, 5})

		out, errCh := Stage(ctx, 2, in, func(ctx context.Context, v int) (int, error) {
			time.Sleep(50 * time.Millisecond)
			return v, nil
		})

		results := []int{}
		errs := []error{}

		// cancel context after 50ms
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		done := make(chan bool)
		go func() {
			for e := range errCh {
				errs = append(errs, e)
			}
			done <- true
		}()

		for v := range out {
			results = append(results, v)
		}
		<-done

		t.Logf("processed %d values before cancellation", len(results))
	})
}

// TestSendOrDone tests the SendOrDone function
func TestSendOrDone(t *testing.T) {
	t.Run("normal send", func(t *testing.T) {
		ctx := context.Background()
		out := make(chan int)

		go func() {
			SendOrDone(ctx, out, 42)
			close(out)
		}()

		v := <-out
		if v != 42 {
			t.Errorf("expect 42, got %d", v)
		}
	})

	t.Run("context canceled - do not send", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		out := make(chan int)

		cancel() // immediately cancel
		SendOrDone(ctx, out, 42)

		select {
		case <-out:
			t.Error("context canceled - should not send value")
		case <-time.After(100 * time.Millisecond):
			// correct, no value sent
		}
	})
}

// Helper function: Mutex for testing
var testMu sync.Mutex
