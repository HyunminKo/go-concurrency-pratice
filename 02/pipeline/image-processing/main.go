package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}

	start := time.Now()

	err := walkFiles(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func walkFiles(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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

		thumbnailImage, err := processImage(path)
		if err != nil {
			return err
		}

		err = saveThumbnail(path, thumbnailImage)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
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

func processImage(path string) (*image.NRGBA, error) {
	srcImage, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}

	thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)
	return thumbnailImage, nil
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
