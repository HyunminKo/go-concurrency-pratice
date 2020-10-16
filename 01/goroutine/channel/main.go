package main

import "fmt"

func main() {
	c := make(chan int)
	go func(a, b int) {
		c <- a + b
	}(1, 2)
	fmt.Printf("computed value %v\n", <-c)

	ch := make(chan int)
	go func() {
		for i := 0; i < 6; i++ {
			ch <- i
		}
		close(ch)
	}()

	for value := range ch {
		fmt.Println(value)
	}

	ch2 := make(chan int, 6)
	go func() {
		defer close(ch2)
		for i := 0; i < 6; i++ {
			fmt.Printf("Sending: %d\n", i)
			ch2 <- i
		}
	}()

	for v := range ch2 {
		fmt.Printf("Received: %d\n", v)
	}
}
