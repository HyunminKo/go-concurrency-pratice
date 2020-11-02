package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
)

func main() {
	req, err := http.NewRequest("GET", "https://andcloud.io", nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1000*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error", err)
		return
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
