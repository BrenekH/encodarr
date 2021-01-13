package controller

import (
	"os"
	"sync"
)

// Queue is a basic implementation of a FIFO Queue.
type Queue struct {
	sync.Mutex
	Items []interface{}
}

// Push appends an item to the end of a Queue.
func (q *Queue) Push(item interface{}) {
	q.Lock()
	defer q.Unlock()
	q.Items = append(q.Items, item)
}

// Pop removes and returns the first item of a Queue.
func (q *Queue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	item := q.Items[0]
	q.Items[0] = nil
	q.Items = q.Items[1:]
	return item
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *Queue) Dequeue() []interface{} {
	q.Lock()
	defer q.Unlock()
	return append(make([]interface{}, 0, len(q.Items)), q.Items...)
}

// IsDirectory returns a bool representing whether or not the provided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
