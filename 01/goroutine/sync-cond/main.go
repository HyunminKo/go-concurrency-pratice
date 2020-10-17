package main

import (
	"fmt"
	"sync"
)

var sharedResource = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	go func() {
		defer wg.Done()

		c.L.Lock()

		for len(sharedResource) == 0 {
			c.Wait()
		}
		fmt.Println(sharedResource["rsc1"])
		c.L.Unlock()
	}()
	c.L.Lock()
	sharedResource["rsc1"] = "foo"
	c.Signal()
	c.L.Unlock()

	wg.Wait()

	wg.Add(2)
	go func() {
		defer wg.Done()

		c.L.Lock()

		for len(sharedResource) < 1 {
			c.Wait()
		}
		fmt.Println(sharedResource["rsc1"])
		c.L.Unlock()
	}()

	go func() {
		defer wg.Done()

		c.L.Lock()

		for len(sharedResource) < 2 {
			c.Wait()
		}
		fmt.Println(sharedResource["rsc2"])
		c.L.Unlock()
	}()

	c.L.Lock()
	sharedResource["rsc1"] = "foo"
	sharedResource["rsc2"] = "bar"
	c.Broadcast()
	c.L.Unlock()

	wg.Wait()
}
