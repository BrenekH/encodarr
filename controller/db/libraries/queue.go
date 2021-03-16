package libraries

import (
	"errors"

	"github.com/BrenekH/project-redcedar-controller/db/dispatched"
)

// Moved from controller/util.go to avoid cyclic imports

// ErrInvalidUUID represents when a passed UUID is invalid
var ErrInvalidUUID error = errors.New("invalid UUID")

// ErrEmptyQueue represents when the operation cannot be completed because the queue is empty
var ErrEmptyQueue error = errors.New("queue is empty")

// Queue is a basic implementation of a FIFO Queue for the Job interface.
type Queue struct {
	Items []dispatched.Job
}

// Push appends an item to the end of a Queue.
func (q *Queue) Push(item dispatched.Job) {
	q.Items = append(q.Items, item)
}

// Pop removes and returns the first item of a Queue.
func (q *Queue) Pop() (dispatched.Job, error) {
	if len(q.Items) == 0 {
		return dispatched.Job{}, ErrEmptyQueue
	}
	item := q.Items[0]
	q.Items[0] = dispatched.Job{} // Hopefully this garbage collects properly
	q.Items = q.Items[1:]
	return item, nil
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *Queue) Dequeue() []dispatched.Job {
	return append(make([]dispatched.Job, 0, len(q.Items)), q.Items...)
}

// InQueue returns a boolean representing whether or not the provided item is in the queue
func (q *Queue) InQueue(item dispatched.Job) bool {
	for _, i := range (*q).Items {
		if item.Equal(i) {
			return true
		}
	}
	return false
}

// InQueuePath returns a boolean representing whether or not the provided item is in the queue based on only the Path field
func (q *Queue) InQueuePath(item dispatched.Job) bool {
	for _, i := range (*q).Items {
		if item.EqualPath(i) {
			return true
		}
	}
	return false
}

// Empty returns a boolean representing whether or not the queue is empty
func (q *Queue) Empty() bool {
	return len(q.Items) == 0
}
