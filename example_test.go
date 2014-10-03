package async_test

import (
	"fmt"
	"time"

	"github.com/egonelbre/async"
)

func ExampleAll_success() {
	result := async.All(
		func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		},
	)

	select {
	case err := <-result.Error:
		fmt.Printf("Got an error: %v\n", err)
	case <-result.Done:
		fmt.Printf("Success\n")
	}
}

func ExampleAll_failing() {
	result := async.All(
		func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(200 * time.Millisecond)
			return fmt.Errorf("CRASH")
		},
	)

	select {
	case err := <-result.Error:
		fmt.Printf("Got an error: %v\n", err)
	case <-result.Done:
		fmt.Printf("Success\n")
	}
}

func ExampleAll_multipleErrors() {
	result := async.All(
		func() error {
			return fmt.Errorf("SPLASH")
		},
		func() error {
			return fmt.Errorf("CRASH")
		},
	)

	for err := range result.Error {
		fmt.Printf("Got an error: %v\n", err)
	}
}

func ExampleSpawn() {
	work := make(chan int, 3)
	done := make(chan int, 3)

	async.Spawn(3, func(id int) {
		for v := range work {
			done <- v * v
		}
	}, func() {
		close(done)
	})

	for i := 0; i < 5; i += 1 {
		work <- i
	}
	close(work)

	for r := range done {
		println(r)
	}
}
