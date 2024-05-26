package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	start := time.Now()
	files := []string{"assets/sample.txt", "assets/sample.jpg", "assets/sample.pdf"}

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			task, err := newTask(file)
			if err != nil {
				return
			}

			task.run()

		}(file)
	}

	wg.Wait()
	log.Println(time.Since(start))
}
