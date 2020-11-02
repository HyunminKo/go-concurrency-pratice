package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

type data struct {
	resutl string
}

func main() {
	deadline := time.Now().Add(10 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	compute := func() <-chan data {
		ch := make(chan data)
		go func() {
			defer close(ch)
			deadline, ok := ctx.Deadline()
			if ok {
				if deadline.Sub(time.Now().Add(50*time.Millisecond)) < 0 {
					fmt.Println("not sufficient time given... terminating")
				}
			}
			time.Sleep(50 * time.Millisecond)

			select {
			case ch <- data{"123"}:
			case <-ctx.Done():
				return
			}
		}()
		return ch
	}

	ch := compute()
	d, ok := <-ch
	if ok {
		fmt.Printf("Work is complete: %s\n", d)
	}
}
