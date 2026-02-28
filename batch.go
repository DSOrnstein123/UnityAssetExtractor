package main

import (
	"fmt"
	"sync"
)

func createBatches(downloadedFiles <-chan string, batches chan<- []string) {
	batch := make([]string, 0, config.BatchSize)

	for filepath := range downloadedFiles {
		batch = append(batch, filepath)

		if len(batch) >= config.BatchSize {
			batches <- batch
			batch = make([]string, 0, config.BatchSize)
		}
	}

	if len(batch) > 0 {
		batches <- batch
	}
}

func batchProcessors(batches <-chan []string) {
	wg := sync.WaitGroup{}

	for i := range config.NumBatchProcessors {
		wg.Add(1)
		go processor(int16(i+1), batches, &wg)
	}

	wg.Wait()
}

func processor(processorId int16, batches <-chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for batch := range batches {
		processBatch(batch, processorId)
	}
}

func processBatch(batch []string, processorId int16) {
	err := decrypt(batch)
	if err != nil {
		fmt.Printf("[BP%d] (Failed) Decrypt batch failed: %v", processorId, err)
	}

	fmt.Printf("[BP%d] ✓ Decrypted %d files\n", processorId, len(batch))

	fmt.Printf("[BP%d] 📦 Extracting %d files...\n", processorId, len(batch))
	err = extract(batch)
	if err != nil {
		fmt.Printf("[BP%d] (Error): %v", processorId, err)
	}

	fmt.Printf("[BP%d] ✅ Batch completed: %d extracted successfully\n\n", processorId, len(batch))
}
