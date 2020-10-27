package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type result struct {
	srcImagePath   string
	thumbnailImage *image.NRGBA
	err            error
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}

	start := time.Now()

	err := setupPipeLine(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func setupPipeLine(root string) error {
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	results := processImage(done, paths)

	for r := range results {
		if r.err != nil {
			return r.err
		}
		saveThumbnail(r.srcImagePath, r.thumbnailImage)
	}

	if err := <-errc; err != nil {
		return err
	}

	return nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {

	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk was canceld")
			}
			return nil
		})
	}()
	return paths, errc
}

func getFileContentType(file string) (string, error) {
	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func processImage(done <-chan struct{}, paths <-chan string) <-chan *result {
	results := make(chan *result)

	thumbnailer := func() {
		for path := range paths {
			srcImage, err := imaging.Open(path)
			if err != nil {
				select {
				case results <- &result{path, nil, err}:
				case <-done:
					return
				}
			}
			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

			select {
			case results <- &result{path, thumbnailImage, nil}:
			case <-done:
				return
			}
		}
	}

	const numThumbnailer = 5
	wg := sync.WaitGroup{}
	wg.Add(numThumbnailer)
	for i := 0; i < numThumbnailer; i++ {
		go func() {
			thumbnailer()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}
