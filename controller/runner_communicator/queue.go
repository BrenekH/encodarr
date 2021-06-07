package runner_communicator

import "github.com/BrenekH/encodarr/controller"

type waitingRunner struct {
	Name         string
	CallbackChan chan controller.Job
}

func newQueue() queue {
	return queue{
		items: make([]waitingRunner, 0),
	}
}

// queue represents a singular queue belonging to one library.
type queue struct {
	items []waitingRunner
}

// Push appends an item to the end of a LibraryQueue.
func (q *queue) Push(item waitingRunner) {
	q.items = append(q.items, item)
}

// Pop removes and returns the first item of a LibraryQueue.
func (q *queue) Pop() (waitingRunner, error) {
	if len(q.items) == 0 {
		return waitingRunner{}, controller.ErrEmptyQueue
	}
	item := q.items[0]
	q.items[0] = waitingRunner{} // Hopefully this garbage collects properly
	q.items = q.items[1:]
	return item, nil
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *queue) Dequeue() []waitingRunner {
	return append(make([]waitingRunner, 0, len(q.items)), q.items...)
}

// Empty returns a boolean representing whether or not the queue is empty
func (q *queue) Empty() bool {
	return len(q.items) == 0
}
