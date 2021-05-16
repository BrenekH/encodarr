package controller

import (
	"reflect"
	"time"
)

type UUID string

// Job represents a job to be carried out by a Runner.
type Job struct {
	UUID     UUID         `json:"uuid"`
	Path     string       `json:"path"`
	Command  []string     `json:"command"`
	Metadata FileMetadata `json:"metadata"`
}

// Library represents a single library.
type Library struct {
	ID                     int           `json:"id"`
	Folder                 string        `json:"folder"`
	Priority               int           `json:"priority"`
	FsCheckInterval        time.Duration `json:"fs_check_interval"`
	Queue                  LibraryQueue  `json:"queue"`
	PathMasks              []string      `json:"path_masks"`
	CommandDeciderSettings string        `json:"command_decider_settings"`
}

type DispatchedJob struct {
	UUID        UUID      `json:"uuid"`
	Runner      string    `json:"runner"`
	Job         Job       `json:"job"`
	Status      JobStatus `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}

type File struct {
	Path     string
	ModTime  time.Time
	Metadata FileMetadata
}

// LibraryQueue represents a singular queue belonging to one library.
type LibraryQueue struct {
	Items []Job
}

// Push appends an item to the end of a LibraryQueue.
func (q *LibraryQueue) Push(item Job) {
	q.Items = append(q.Items, item)
}

// Pop removes and returns the first item of a LibraryQueue.
func (q *LibraryQueue) Pop() (Job, error) {
	if len(q.Items) == 0 {
		return Job{}, ErrEmptyQueue
	}
	item := q.Items[0]
	q.Items[0] = Job{} // Hopefully this garbage collects properly
	q.Items = q.Items[1:]
	return item, nil
}

// Dequeue returns a copy of the underlying slice in the Queue.
func (q *LibraryQueue) Dequeue() []Job {
	return append(make([]Job, 0, len(q.Items)), q.Items...)
}

// InQueue returns a boolean representing whether or not the provided item is in the queue
func (q *LibraryQueue) InQueue(item Job) bool {
	for _, i := range (*q).Items {
		if item.Equal(i) {
			return true
		}
	}
	return false
}

// InQueuePath returns a boolean representing whether or not the provided item is in the queue based on only the Path field
func (q *LibraryQueue) InQueuePath(item Job) bool {
	for _, i := range (*q).Items {
		if item.EqualPath(i) {
			return true
		}
	}
	return false
}

// Empty returns a boolean representing whether or not the queue is empty
func (q *LibraryQueue) Empty() bool {
	return len(q.Items) == 0
}

// Equal is a custom equality check for the Job type
func (j Job) Equal(check Job) bool {
	if j.UUID != check.UUID {
		return false
	}
	if j.Path != check.Path {
		return false
	}
	if !reflect.DeepEqual(j.Command, check.Command) {
		return false
	}
	return true
}

// EqualPath is a custom equality check for the Job type that only checks the Path parameter
func (j Job) EqualPath(check Job) bool {
	return j.Path == check.Path
}
