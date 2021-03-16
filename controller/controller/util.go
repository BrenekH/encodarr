package controller

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/BrenekH/project-redcedar-controller/db/dispatched"
)

// ErrInvalidUUID represents when a passed UUID is invalid
var ErrInvalidUUID error = errors.New("invalid UUID")

// ErrEmptyQueue represents when the operation cannot be completed because the queue is empty
var ErrEmptyQueue error = errors.New("queue is empty")

// Queue is a basic implementation of a FIFO Queue for the Job interface.
type Queue struct {
	sync.Mutex
	items []dispatched.Job
}

// Push appends an item to the end of a Queue.
func (q *Queue) Push(item dispatched.Job) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

// Pop removes and returns the first item of a Queue.
func (q *Queue) Pop() (dispatched.Job, error) {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return dispatched.Job{}, ErrEmptyQueue
	}
	item := q.items[0]
	q.items[0] = dispatched.Job{} // Hopefully this garbage collects properly
	q.items = q.items[1:]
	return item, nil
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *Queue) Dequeue() []dispatched.Job {
	q.Lock()
	defer q.Unlock()
	return append(make([]dispatched.Job, 0, len(q.items)), q.items...)
}

// InQueue returns a boolean representing whether or not the provided item is in the queue
func (q *Queue) InQueue(item dispatched.Job) bool {
	q.Lock()
	defer q.Unlock()
	for _, i := range (*q).items {
		if item.Equal(i) {
			return true
		}
	}
	return false
}

// InQueuePath returns a boolean representing whether or not the provided item is in the queue based on only the Path field
func (q *Queue) InQueuePath(item dispatched.Job) bool {
	q.Lock()
	defer q.Unlock()
	for _, i := range (*q).items {
		if item.EqualPath(i) {
			return true
		}
	}
	return false
}

// Empty returns a boolean representing whether or not the queue is empty
func (q *Queue) Empty() bool {
	q.Lock()
	defer q.Unlock()
	return len(q.items) == 0
}

// IsDirectory returns a bool representing whether or not the provided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// MoveFile moves the sourcePath to the destPath without tripping up on file system issues
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	return nil
}
