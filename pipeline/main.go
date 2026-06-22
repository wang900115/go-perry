package main

import (
	"context"
	"errors"
	"sync"

	"github.com/labstack/gommon/log"
	"golang.org/x/sync/semaphore"
)

// ReadOrDone is a generic function that reads from an input channel and returns a new channel that will close when the context is done or when the input channel is closed.
func ReadOrDone[T any](ctx context.Context, in <-chan T) <-chan T {
	valStream := make(chan T)

	go func() {
		defer close(valStream)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-ctx.Done():
				}
			}
		}
	}()
	return valStream
}

// SendOrDone is a generic function that sends a value to an output channel unless the context is done, in which case it returns without sending.
func SendOrDone[T any](ctx context.Context, out chan<- T, v T) {
	select {
	case <-ctx.Done():
		return
	case out <- v:
	}
}

// FromSlice is a generic function that creates a channel from a slice of values.
func FromSlice[T any](ctx context.Context, slice []T) <-chan T {
	out := make(chan T, len(slice))
	for _, v := range slice {
		out <- v
	}
	close(out)
	return out
}

// Merge is a generic function that merges multiple input channels into a single output channel. It reads from each input channel and sends the values to the output channel until all input channels are closed or the context is done.
func Merge[T any](ctx context.Context, cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()
		for v := range ReadOrDone(ctx, c) {
			select {
			case out <- v:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Tee is a generic function that duplicates the values from an input channel into two output channels. It reads from the input channel and sends the values to both output channels until the input channel is closed or the context is done.
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	// Use buffered channels to avoid deadlock when consumption speeds differ
	out1 := make(chan T, 100)
	out2 := make(chan T, 100)

	go func() {
		defer close(out1)
		defer close(out2)

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}

				// Send to both channels (buffered, won't block)
				out1 <- val
				out2 <- val
			}
		}
	}()

	return out1, out2
}

// Stage is a generic function that processes values from an input channel using a provided function and sends the results to an output channel. It uses a semaphore to limit the number of concurrent workers processing the input values. The function will return two channels: one for the processed output and one for any errors that occur during processing.
func Stage[In any, Out any](ctx context.Context, maxWorker int, in <-chan In, fn func(context.Context, In) (Out, error)) (chan Out, chan error) {
	outputChannl := make(chan Out)
	errorChannel := make(chan error)

	limit := int64(maxWorker)

	sem1 := semaphore.NewWeighted(limit)

	go func() {
		defer close(outputChannl)
		defer close(errorChannel)

		for s := range ReadOrDone(ctx, in) {
			if err := sem1.Acquire(ctx, 1); err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Errorf("failed to acquire semaphore: %v", err)
				}
				break
			}

			go func(s In) {
				defer sem1.Release(1)

				result, err := fn(ctx, s)
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						errorChannel <- err
					}
				} else {
					outputChannl <- result
				}
			}(s)
		}

		if err := sem1.Acquire(context.Background(), limit); err != nil {
			log.Errorf("failed to acquire semaphore for shutdown: %v", err)
		}
	}()

	return outputChannl, errorChannel
}
