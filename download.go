package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

func downloadAllFiles(urls []string, output string, outputChan chan<- string) {
	jobs := make(chan string, totalFiles)
	wg := sync.WaitGroup{}

	for i := range config.MaxDownloadWorkers {
		go downloadWorker(int16(i+1), jobs, output, outputChan, &wg)
	}

	for _, url := range urls {
		jobs <- url
		wg.Add(1)
	}

	close(jobs)
	wg.Wait()
}

func downloadWorker(workderId int16, jobs <-chan string, output string, outputChan chan<- string, wg *sync.WaitGroup) {
	for url := range jobs {
		filename, err := downloadFile(url, output)
		if err != nil {
			fmt.Printf("[DW%d] (Failed) %s - %v\n", workderId, filename, err)
			wg.Done()
			continue
		}

		filepath := filepath.Join(output, filename)
		outputChan <- filepath

		current := atomic.AddInt32(&downloadCompleted, 1)
		fmt.Printf("[DW%d] %s (%d/%d)\n", workderId, filename, current, totalFiles)
		wg.Done()
	}
}

func downloadFile(url string, output string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	filename := filepath.Base(resp.Request.URL.Path)
	outputPath := filepath.Join(output, filename)

	out, err := os.Create(outputPath)
	if err != nil {
		return "", nil
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", nil
	}

	return filename, nil
}
