package controller

import (
	"errors"
	"os"
	"sync"
)

// ErrInvalidUUID represents when a passed UUID is invalid
var ErrInvalidUUID error = errors.New("Invalid UUID")

// ErrEmptyQueue represents when the operation cannot be completed because the queue is empty
var ErrEmptyQueue error = errors.New("Queue is empty")

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
		return Job{}, ErrEmptyQueue
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

// InContainerPath returns a boolean representing whether or not the provided Job is in the container based on only the Path field
func (c *DispatchedContainer) InContainerPath(item Job) bool {
	c.Lock()
	defer c.Unlock()
	for _, v := range (*c).items {
		if item.EqualPath(v.Job) {
			return true
		}
	}
	return false
}

// UpdateStatus uses the provided UUID string to identify the Job to be updated with the new status as defined by the provided JobStatus
func (c *DispatchedContainer) UpdateStatus(u string, js JobStatus) error {
	c.Lock()
	defer c.Unlock()
	for i, v := range c.items {
		if v.Job.EqualUUID(Job{UUID: u}) {
			// Save before removing from container slice
			interim := v

			// Remove from container slice
			c.items[i] = c.items[len(c.items)-1]
			c.items[len(c.items)-1] = DispatchedJob{}
			c.items = c.items[:len(c.items)-1]

			// Add back into container slice with modifications
			interim.Status = js
			c.items = append(c.items, interim)
			return nil
		}
	}
	return ErrInvalidUUID
}

// IsDirectory returns a bool representing whether or not the provided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
