package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Result struct {
	URL      string
	Filename string
	Err      error
}
func downloadFile(id int, url string, semaphore chan struct{}, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	semaphore <- struct{}{}
	defer func() {
		<-semaphore
	}()

	// Give every file a unique name
	filename := fmt.Sprintf("image_%d.jpg", id)
	savePath := filepath.Join("downloads", filename)

	fmt.Printf("Starting download: %s\n", filename)

	resp, err := http.Get(url)
	if err != nil {
		results <- Result{
			URL:      url,
			Filename: filename,
			Err:      err,
		}
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(savePath)
	if err != nil {
		results <- Result{
			URL:      url,
			Filename: filename,
			Err:      err,
		}
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		results <- Result{
			URL:      url,
			Filename: filename,
			Err:      err,
		}
		return
	}

	results <- Result{
		URL:      url,
		Filename: savePath,
		Err:      nil,
	}
}

func main() {
	// Create downloads directory if it doesn't exist
	err := os.MkdirAll("downloads", os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create downloads folder: %v\n", err)
		return
	}

	urls := []string{
		"https://picsum.photos/200/300",
		"https://picsum.photos/300/300",
		"https://picsum.photos/400/300",
		"https://picsum.photos/500/300",
		"https://picsum.photos/600/300",
	}

	semaphore := make(chan struct{}, 3) // only 3 concurrent downloads
	results := make(chan Result)

	 wg := sync.WaitGroup{}
	 
	for i, url := range urls {
		wg.Add(1)
		go downloadFile(i+1, url, semaphore, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err != nil {
			fmt.Printf("Failed: %s -> %v\n", result.URL, result.Err)
			continue
		}

		fmt.Printf("Finished downloading: %s\n", result.Filename)
	}

	fmt.Println("All downloads completed.")
}