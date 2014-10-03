package async

import "sync/atomic"

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
