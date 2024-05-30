package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var files = [4]string{"assets/sample.txt", "assets/sample.jpg", "assets/sample.pdf", "assets/sample.mkv"}

func main() {
	var wg sync.WaitGroup
	start := time.Now()

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			task, err := newTask(file)
			if err != nil {
				log.Println("unable to find the target file", file)
				return
			}

			err = task.run()
			if err != nil {
				log.Println("unable to process", file)
				return
			}

			os.Remove(task.getOutFileName())
		}(file)
	}

	wg.Wait()
	log.Println(time.Since(start))
}
