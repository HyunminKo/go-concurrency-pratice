package main

import "fmt"

func genMsg(ch1 chan<- int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("Sending: %d\n", i)
		ch1 <- i
	}
	close(ch1)
}

func relayMsg(ch1 <-chan int, ch2 chan<- int) {
	for v := range ch1 {
		fmt.Printf("Relaying: %d\n", v)
		ch2 <- v
	}
	close(ch2)
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go genMsg(ch1)
	go relayMsg(ch1, ch2)

	for v := range ch2 {
		fmt.Printf("Received: %d\n", v)
	}
}
