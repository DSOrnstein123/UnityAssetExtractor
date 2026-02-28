package main

import (
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

	links := []string{}

	totalFiles = int32(len(links))

	// gui(links)
	processWithBatching(links)
}

func processWithBatching(urls []string) {
	downloadedFiles := make(chan string, config.BatchSize*2)
	batches := make(chan []string, config.NumBatchProcessors)

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
