package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("Allocated new bytes.Buffer")
		return new(bytes.Buffer)
	},
}

func log(w io.Writer, debug string) {

	b := bufPool.Get().(*bytes.Buffer)

	b.Reset()

	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(debug)
	b.WriteString("\n")
	w.Write(b.Bytes())

	bufPool.Put(b)
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			defer wg.Done()
			log(os.Stdout, fmt.Sprintf("debug-string: %d", i))
		}(i)
	}
	wg.Wait()

	log(os.Stdout, fmt.Sprintf("debug-string: %d", 0))
	log(os.Stdout, fmt.Sprintf("debug-string: %d", 1))
}
