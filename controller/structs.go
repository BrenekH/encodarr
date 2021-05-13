package controller

import "time"

type UUID string

// Job represents a job to be carried out by a Runner.
type Job struct{}

// LibrarySettings represents the settings of a single library.
type LibrarySettings struct{}

// LibraryQueue represents a singular queue belonging to one library.
type LibraryQueue struct{}

type DispatchedJob struct {
	UUID        UUID
	Runner      string
	Job         Job
	Status      JobStatus
	LastUpdated time.Time
}

type JobStatus struct{}
