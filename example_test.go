package async_test

import (
	"fmt"
	"sync/atomic"
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
		fmt.Println(r)
	}
}

func ExampleRun() {
	total := int64(0)

	async.Run(8, func(id int) {
		atomic.AddInt64(&total, int64(id))
	})

	fmt.Println(total)
	// Output:28
}

func ExampleIter() {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	output := make([]int, len(input))
	async.Iter(len(input), 2, func(i int) {
		output[i] = input[i] * input[i]
	})

	for i, v := range output {
		fmt.Println(i, input[i], v)
	}
}

func ExampleBlockIter() {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	output := make([]int, len(input))
	async.BlockIter(len(input), 2, func(start, limit int) {
		for i, v := range input[start:limit] {
			output[start+i] = v * v
		}
	})

	for i, v := range output {
		fmt.Println(i, input[i], v)
	}
}
