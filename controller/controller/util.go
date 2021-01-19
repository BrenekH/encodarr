package controller

import (
	"os"
	"sync"
)

// Queue is a basic implementation of a FIFO Queue.
type Queue struct {
	sync.Mutex
	items []interface{}
}

// Push appends an item to the end of a Queue.
func (q *Queue) Push(item interface{}) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

// Pop removes and returns the first item of a Queue.
func (q *Queue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items[0] = nil
	q.items = q.items[1:]
	return item
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *Queue) Dequeue() []interface{} {
	q.Lock()
	defer q.Unlock()
	return append(make([]interface{}, 0, len(q.items)), q.items...)
}

// InQueue returns a boolean representing whether or not the provided item is in the queue
func (q *Queue) InQueue(item interface{}) bool {
	q.Lock()
	defer q.Unlock()
	for _, i := range (*q).items {
		if i == item {
			return true
		}
	}
	return false
}

// IsDirectory returns a bool representing whether or not the provided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
