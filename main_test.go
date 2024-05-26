package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func BenchmarkTransfer(b *testing.B) {
	var wg sync.WaitGroup

	files := []string{"sample.txt", "sample.jpg", "sample.pdf"}

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			task, err := newTask(file)
			if err != nil {
				return
			}

			task.run()
			fmt.Println("========================")

			os.Remove(task.getOutFileName())
		}(file)
	}

	wg.Wait()
}
