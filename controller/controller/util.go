package controller

import (
	"fmt"
	"os"
	"sync"
)

// Queue is a basic implementation of a FIFO Queue for the Job interface.
type Queue struct {
	sync.Mutex
	items []Job
}

// Push appends an item to the end of a Queue.
func (q *Queue) Push(item Job) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

// Pop removes and returns the first item of a Queue.
func (q *Queue) Pop() (Job, error) {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return Job{}, fmt.Errorf("Queue is empty")
	}
	item := q.items[0]
	q.items[0] = Job{} // Hopefully this garbage collects properly
	q.items = q.items[1:]
	return item, nil
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *Queue) Dequeue() []Job {
	q.Lock()
	defer q.Unlock()
	return append(make([]Job, 0, len(q.items)), q.items...)
}

// InQueue returns a boolean representing whether or not the provided item is in the queue
func (q *Queue) InQueue(item Job) bool {
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
func (q *Queue) InQueuePath(item Job) bool {
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

// DispatchedContainer is a container struct for dispatched jobs
type DispatchedContainer struct {
	sync.Mutex
	items []DispatchedJob
}

// Add adds the supplied DispatchedJob to the container
func (c *DispatchedContainer) Add(item DispatchedJob) {
	c.Lock()
	defer c.Unlock()

	c.items = append(c.items, item)
}

// Decontain returns a copy of the underlying slice in the Container.
func (c *DispatchedContainer) Decontain() []DispatchedJob {
	c.Lock()
	defer c.Unlock()
	return append(make([]DispatchedJob, 0, len(c.items)), c.items...)
}

// IsDirectory returns a bool representing whether or not the provided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
