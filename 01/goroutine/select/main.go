package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			fmt.Printf("Received from ch1: %v\n", v)
		case v := <-ch2:
			fmt.Printf("Received from ch2: %v\n", v)
		}
	}

	ch3 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		ch3 <- "three"
	}()

	select {
	case v := <-ch3:
		fmt.Println(v)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout!!")
	}

	ch4 := make(chan string)
	go func() {
		for i := 0; i < 3; i++ {
			ch4 <- "message"
		}
	}()

	for i := 0; i < 2; i++ {
		select {
		case v := <-ch4:
			fmt.Println(v)
		default:
			fmt.Println("No message")
		}
		fmt.Println("processing")
		time.Sleep(1500 * time.Millisecond)
	}

}
