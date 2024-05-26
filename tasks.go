package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const splitInto int64 = 5

type TaskResult struct {
	chunk []byte
	part  int
}

type WorkerPool struct {
	tasks      []*Task
	bufferSize uint32
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		tasks:      make([]*Task, 0),
		bufferSize: 5,
	}
}

type Task struct {
	filePath string
	size     int64
}

func newTask(path string) (*Task, error) {
	task := &Task{
		filePath: path,
	}
	size, err := task.getFileSize()
	if err != nil {
		log.Println("invalid file path", path)
		return nil, err
	}

	task.size = size
	return task, nil
}

func (t *Task) getFileSize() (int64, error) {
	f, err := os.Stat(t.filePath)
	if err != nil {
		log.Println("unable to get file details", t.filePath)
		return 0, err
	}

	return f.Size(), nil
}

func (t *Task) splitFile() error {
	fileSize, err := t.getFileSize()
	if err != nil {
		return err
	}

	sizePerPart := (fileSize + splitInto - 1) / splitInto
	log.Println("Size per part:", sizePerPart)

	writeChannel := make(chan TaskResult, splitInto)

	var wg sync.WaitGroup

	for part, end := 0, int64(0); end < fileSize; part, end = part+1, end+sizePerPart {
		wg.Add(1)

		start := end
		nextEnd := end + sizePerPart
		if nextEnd > fileSize {
			nextEnd = fileSize
		}

		go func(part int, start, sizePerPart int64) {
			buffer, err := t.readChunk(start, sizePerPart)
			if err != nil {
				log.Fatalln("unable to read chunk:", err)
			}

			writeChannel <- TaskResult{
				chunk: buffer,
				part:  part,
			}

			wg.Done()
		}(part, start, sizePerPart)
	}

	wg.Wait()
	close(writeChannel)

	result := map[int][]byte{}
	var mu sync.RWMutex

	wg.Add(1)
	go func(result map[int][]byte, mu *sync.RWMutex) {
		defer wg.Done()
		for taskResult := range writeChannel {
			result[taskResult.part] = taskResult.chunk
		}
	}(result, &mu)

	wg.Wait()

	outfile := t.getOutFileName()

	for i := 0; i < int(splitInto); i++ {
		err = t.writeChunk(outfile, result[i], int(splitInto))
		if err != nil {
			log.Fatalln("unable to write to outfile", err)
			return err
		}
	}

	return nil
}

func (t *Task) getOutFileName() string {
	ext := filepath.Ext(t.filePath)
	fileName := filepath.Base(t.filePath)
	pos := strings.LastIndexByte(fileName, '.')
	if pos == -1 {
		log.Fatalln("unable to get outfile name")
	}

	name := fileName[:pos]
	outfile := fmt.Sprintf("out/%s_out%s", name, ext)
	return outfile
}

func (t *Task) readChunk(start, sizePerPart int64) ([]byte, error) {
	file, err := os.Open(t.filePath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	buffer := make([]byte, sizePerPart)
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buffer[:n], nil
}

func (t *Task) writeChunk(outfile string, buffer []byte, _ int) error {
	file, err := os.OpenFile(outfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("unable to open outfile for writing", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer)
	if err != nil {
		log.Println("unable to write to output file", err)
		return err
	}

	return nil
}

func (t *Task) checkIntegrity() error {
	in, err := getMD5sum(t.filePath)
	if err != nil {
		return err
	}

	out, err := getMD5sum(t.getOutFileName())
	if err != nil {
		return err
	}

	if out != in {
		return errors.New("checksums of input and output do not match")
	}

	return nil
}

func (t *Task) run() error {
	size, err := t.getFileSize()
	if err != nil {
		log.Println("unable to file size")
		return err
	}

	log.Printf("Target file: %s\n", t.filePath)
	log.Printf("File size: %vkb\n", size)

	err = t.splitFile()
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.checkIntegrity()
	if err != nil {
		log.Println("checksum do not match")
		return err
	}

	log.Println("File transfer completed!")
	return nil
}

func getMD5sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
