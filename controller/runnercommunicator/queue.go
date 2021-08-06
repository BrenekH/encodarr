package runnercommunicator

import (
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

type waitingRunner struct {
	Name         string
	CallbackChan chan controller.Job
	UUID         string
}

func newQueue() queue {
	return queue{
		items: make([]waitingRunner, 0),
	}
}

// queue represents a singular queue belonging to one library.
type queue struct {
	sync.Mutex
	items []waitingRunner
}

// Push appends an item to the end of a LibraryQueue.
func (q *queue) Push(item waitingRunner) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

// Pop removes and returns the first item of a LibraryQueue.
func (q *queue) Pop() (waitingRunner, error) {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return waitingRunner{}, controller.ErrEmptyQueue
	}
	item := q.items[0]
	q.items[0] = waitingRunner{} // Hopefully this garbage collects properly
	q.items = q.items[1:]
	return item, nil
}

// Remove deletes the first item that has the uuid provided.
func (q *queue) Remove(uuid string) {
	q.Lock()
	defer q.Unlock()
	for index, v := range q.items {
		if v.UUID == uuid {
			q.items = append(q.items[:index], q.items[index+1:]...)
			return
		}
	}
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *queue) Dequeue() []waitingRunner {
	q.Lock()
	defer q.Unlock()
	return append(make([]waitingRunner, 0, len(q.items)), q.items...)
}

// Empty returns a boolean representing whether or not the queue is empty
func (q *queue) Empty() bool {
	q.Lock()
	defer q.Unlock()
	return len(q.items) == 0
}
