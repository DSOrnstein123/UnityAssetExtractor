package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	RootURL            string
	MaxDownloadWorkers int
	NumBatchProcessors int
	BatchSize          int
	DownloadFolder     string
	DecryptScript      string
}

var (
	config Config

	downloadCompleted int32
	totalFiles        int32
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config = Config{
		RootURL:            os.Getenv("ROOT_URL"),
		MaxDownloadWorkers: mustInt("MAX_DOWNLOAD_WORKERS"),
		NumBatchProcessors: mustInt("NUM_BATCH_PROCESSORS"),
		BatchSize:          mustInt("BATCH_SIZE"),
		DownloadFolder:     os.Getenv("DOWNLOAD_FOLDER"),
		DecryptScript:      os.Getenv("DECRYPT_SCRIPT"),
	}

	gui()
}

func createUrls() ([]string, error) {
	data, err := os.ReadFile("list.json")
	if err != nil {
		fmt.Printf("Lỗi khi đọc file: %v\n", err)
		return nil, err
	}

	var assetList []string
	err = json.Unmarshal(data, &assetList)
	if err != nil {
		fmt.Printf("Lỗi khi giải mã JSON: %v\n", err)
		return nil, err
	}

	var urls []string
	for _, elem := range assetList {
		urls = append(urls, config.RootURL+elem)
	}

	return urls, nil
}

func processWithBatching() {
	downloadedFiles := make(chan string, config.BatchSize*2)
	batches := make(chan []string, config.NumBatchProcessors)

	urls, err := createUrls()
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		defer close(downloadedFiles)
		downloadAllFiles(urls, config.DownloadFolder, downloadedFiles)
	}()

	go func() {
		defer wg.Done()
		defer close(batches)
		createBatches(downloadedFiles, batches)
	}()

	go func() {
		defer wg.Done()
		batchProcessors(batches)
	}()

	wg.Wait()
}
