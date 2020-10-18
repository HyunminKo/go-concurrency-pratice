package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/html"
)

var fetched map[string]bool

type result struct {
	url   string
	urls  []string
	err   error
	depth int
}

func Crawl(url string, depth int) {
	if depth < 0 {
		return
	}

	urls, err := findLinks(url)

	if err != nil {
		return
	}

	fmt.Printf("found: %s\n", url)
	fetched[url] = true
	for _, u := range urls {
		if !fetched[u] {
			Crawl(u, depth-1)
		}
	}
	return
}

func ConcurrentCrawl(url string, depth int) {
	ch := make(chan *result)

	fetch := func(url string, depth int) {
		urls, err := findLinks(url)
		ch <- &result{url, urls, err, depth}
	}

	go fetch(url, depth)
	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		res := <-ch
		if res.err != nil {
			continue
		}
		fmt.Printf("found: %s\n", res.url)

		if res.depth > 0 {
			for _, u := range res.urls {
				if !fetched[u] {
					fetching++
					go fetch(u, res.depth-1)
					fetched[u] = true
				}
			}
		}
	}
	close(ch)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fetched = make(map[string]bool)
	now := time.Now()
	ConcurrentCrawl("http://andcloud.io", 2)
	fmt.Println("time taken:", time.Since(now))
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTMl: %v", url, err)
	}
	return visit(nil, doc), nil
}
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
