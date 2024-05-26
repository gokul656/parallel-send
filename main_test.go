package main

import (
	"os"
	"sync"
	"testing"
)

func BenchmarkTransfer(b *testing.B) {
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			task, err := newTask(file)
			if err != nil {
				return
			}

			task.run()

			os.Remove(task.getOutFileName())
		}(file)
	}

	wg.Wait()
}
