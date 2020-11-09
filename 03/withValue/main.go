package main

import (
	"context"
	"fmt"
)

type database map[string]bool
type userIDKeyType string 
var db database = database{
	"jane": true,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	processRequest(ctx,"jane")
}

func processRequest(ctx context.Context, userID string) {
	vctx := context.WithValue(ctx, userIDKeyType("userIDKey"),"jane")
	ch := checkMembership(vctx)
	status := <-ch
	fmt.Printf("membership status of userId : %s : %v\n", userID, status)
}

func checkMembership(ctx context.Context) <-chan bool {
	ch := make(chan bool)
	go func() {
		defer close(ch)
		userID := ctx.Value(userIDKeyType("userIDKey")).(string)
		status := db[userID]
		ch <- status
	}()
	return ch
}