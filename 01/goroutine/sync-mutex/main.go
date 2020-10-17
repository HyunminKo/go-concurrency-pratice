package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var balance int
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	deposit := func(amount int) {
		mu.Lock()
		defer mu.Unlock()
		balance += amount
	}

	withdrawal := func(amount int) {
		mu.Lock()
		defer mu.Unlock()
		balance -= amount
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			deposit(1)
		}()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			withdrawal(1)
		}()
	}
	wg.Wait()
	fmt.Println(balance)
}
