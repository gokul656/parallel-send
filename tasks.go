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
)

const splitInto int64 = 5

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
		log.Println("invalid file path")
		return nil, err
	}

	task.size = size
	return task, nil
}

func (t *Task) getFileSize() (int64, error) {
	f, err := os.Stat(t.filePath)
	if err != nil {
		return 0, err
	}

	return f.Size(), nil
}

func (t *Task) splitFile() error {
	file, err := os.Open(t.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileSize, err := t.getFileSize()
	if err != nil {
		return err
	}

	sizePerPart := (fileSize + splitInto - 1) / splitInto
	log.Println("Size per part:", sizePerPart)

	for end := int64(0); end < fileSize; end += sizePerPart {
		start := end
		if end != 0 {
			start++
		}
		if start >= fileSize {
			break
		}

		nextEnd := end + sizePerPart
		if nextEnd > fileSize {
			nextEnd = fileSize
		}

		buffer, err := t.readChunk(file, start, nextEnd-start)
		if err != nil {
			log.Fatalln("unable to read chunk:", err)
		}

		outfile := t.getOutFileName()
		err = t.writeChunk(outfile, buffer)
		if err != nil {
			log.Fatalln("unable to write to outfile:", err)
		}
	}

	return t.checkIntegrity()
}

func (t *Task) getOutFileName() string {
	ext := filepath.Ext(t.filePath)
	fileName := filepath.Base(t.filePath)
	pos := strings.LastIndexByte(fileName, '.')
	if pos == -1 {
		log.Fatalln("unable to get outfile name")
	}

	name := fileName[:pos]
	outfile := fmt.Sprintf("%s_out%s", name, ext)
	return outfile
}

func (t *Task) readChunk(file *os.File, start, sizePerPart int64) ([]byte, error) {
	buffer := make([]byte, sizePerPart+1)
	_, err := file.Seek(start, io.SeekStart)
	if err != nil {
		return nil, err
	}

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buffer[:n], nil
}

func (t *Task) writeChunk(outfile string, buffer []byte) error {
	file, err := os.OpenFile(outfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("unable to open outfile for writing:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer)
	return err
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

func (t *Task) run() {
	size, err := t.getFileSize()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Target file: %s\n", t.filePath)
	log.Printf("File size: %vkb\n", size)

	err = t.splitFile()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("File transfer completed!")
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
