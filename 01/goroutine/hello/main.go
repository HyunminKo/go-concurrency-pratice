package main

import (
	"fmt"
	"time"
)

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	fun("direct call")

	go fun("goroutine-1")

	go func() {
		fun("goroutine-2")
	}()

	fv := fun
	go fv("goroutine-3")

	fmt.Println("wait for goroutines")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("done...")
}
