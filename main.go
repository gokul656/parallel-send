package main

import (
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	start := time.Now()
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

			os.Remove(task.getOutFileName())
		}(file)
	}

	wg.Wait()
	log.Println(time.Since(start))
}
