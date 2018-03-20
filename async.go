package async

import (
	"sync"
	"sync/atomic"
)

type Result struct {
	Done  <-chan struct{}
	Error <-chan error
}

// All starts all functions concurrently
// if any error occurs it will be sent to the Error channel
// after all functions have terminated the Done channel will get a single value
func All(fns ...func() error) Result {
	done := make(chan struct{}, 1)
	errs := make(chan error, len(fns))

	waiting := int32(len(fns))

	for _, fn := range fns {
		go func(fn func() error) {
			if err := fn(); err != nil {
				errs <- err
			}

			if atomic.AddInt32(&waiting, -1) == 0 {
				done <- struct{}{}
				close(errs)
				close(done)
			}
		}(fn)
	}

	if len(fns) == 0 {
		done <- struct{}{}
		close(errs)
		close(done)
	}

	return Result{done, errs}
}

// Spawns N routines, after each completes runs all whendone functions
func Spawn(N int, fn func(id int), whendone ...func()) {
	waiting := int32(N)
	for k := 0; k < N; k += 1 {
		go func(k int) {
			fn(k)
			if atomic.AddInt32(&waiting, -1) == 0 {
				for _, fn := range whendone {
					fn()
				}
			}
		}(int(k))
	}
}

// Run N routines and wait for all to complete
func Run(N int, fn func(id int)) {
	var wg sync.WaitGroup
	wg.Add(N)
	for k := 0; k < N; k += 1 {
		go func(k int) {
			fn(k)
			wg.Done()
		}(int(k))
	}
	wg.Wait()
}

// Spawns N routines, iterating over [0..Count) items in increasing order
func Iter(Count int, N int, fn func(i int)) {
	var wg sync.WaitGroup
	wg.Add(N)
	i := int64(0)
	for k := 0; k < N; k += 1 {
		go func() {
			defer wg.Done()
			for {
				idx := int(atomic.AddInt64(&i, 1) - 1)
				if idx >= Count {
					break
				}
				fn(idx)
			}
		}()
	}
	wg.Wait()
}

// Spawns N routines, iterating over [0..Count] items by splitting
// them into blocks [start..limit), note that item "limit" shouldn't be
// processed.
func BlockIter(Count int, N int, fn func(start, limit int)) {
	var wg sync.WaitGroup

	start, left := 0, Count
	for k := 0; k < N; k += 1 {
		count := (left + (N - k - 1)) / (N - k)
		limit := start + count
		if limit >= Count {
			limit = Count
		}
		wg.Add(1)
		go func(start, limit int) {
			defer wg.Done()
			fn(start, limit)
		}(start, limit)
		start = start + count
		left -= count
		if left <= 0 {
			break
		}
	}
	wg.Wait()
}
