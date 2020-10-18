package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func generator(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- v * v:
			case <-done:
				return
			}
		}
	}()
	return out
}

// fan in
func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			select {
			case out <- n:
			case <-done:
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

func main() {
	done := make(chan struct{})
	in := generator(done, 2, 3)

	ch1 := square(done, in)
	ch2 := square(done, in)

	out := merge(done, ch1, ch2)

	fmt.Println(<-out)
	close(done)

	time.Sleep(10 * time.Millisecond)

	g := runtime.NumGoroutine()
	fmt.Printf("number of goroutines active = %d\n", g)
}
