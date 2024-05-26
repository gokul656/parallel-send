package main

import (
	"log"
	"sync"
	"time"
)

var files = [3]string{"assets/sample.txt", "assets/sample.jpg", "assets/sample.pdf"}

func main() {
	var wg sync.WaitGroup

	start := time.Now()

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
