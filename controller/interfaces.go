package controller

import (
	"context"
	"sync"
	"time"
)

// The HealthChecker interface describes how a struct wishing to decide if a job's
// last update was long enough ago to mark the Runner doing it as unresponsive
// should interact with the Run function.
type HealthChecker interface {
	// Run loops through the provided slice of dispatched jobs and checks if any have
	// surpassed the allowed time between updates.
	Run() (uuidsToNull []UUID)

	Start(ctx *context.Context)
}

// The LibraryManager interface describes how a struct wishing to deal with user's
// libraries should interact with the Run function.
type LibraryManager interface {
	// ImportCompletedJobs imports the provided jobs into the system.
	ImportCompletedJobs([]Job)

	// LibrarySettings returns the current settings of all libraries
	LibrarySettings() []Library

	// LibraryQueues returns a slice of LibraryQueue to indicate the current status of
	// the queues.
	LibraryQueues() []LibraryQueue

	// PopNewJob returns a job that may be dispatched as well as deletes it from any
	// data stores.
	PopNewJob() Job

	// UpdateLibrarySettings loops through the provided map of new settings and applies
	// them to the appropriate libraries.
	UpdateLibrarySettings(map[string]Library)

	Start(ctx *context.Context, wg *sync.WaitGroup)
}

// The RunnerCommunicator interface describes how a struct wishing to communicate
// with external Runners should interact with the Run function.
type RunnerCommunicator interface {
	// CompletedJobs returns a slice of jobs that are ready to be imported back into the
	// system.
	CompletedJobs() []Job

	// NewJob takes the provided job and sends it to a waiting Runner.
	NewJob(Job)

	// NeedNewJob returns a boolean indicating whether or not a new job is required.
	NeedNewJob() bool

	// NullifyUUIDs takes the provided slice of UUIDs and marks them
	// so that if a Runner sends a request with a nullified UUID, it gets notified
	// that it is considered unresponsive and should acquire a new job.
	NullifyUUIDs([]UUID)

	// WaitingRunners returns the names of all the Runners which are waiting for a job.
	WaitingRunners() (runnerNames []string)

	Start(ctx *context.Context, wg *sync.WaitGroup)
}

// The UserInterfacer interface describes how a struct wishing to interact
// with the user should interact with the Run function.
type UserInterfacer interface {
	// NewLibrarySettings returns a map of all updated library settings as set by the user.
	NewLibrarySettings() map[string]Library

	// SetLibrarySettings takes the provided slice of LibrarySettings and stores it
	// for an incoming request.
	SetLibrarySettings([]Library)

	// SetLibraryQueues takes the provided slice and stores it so that if a request
	// to view the queues is received, the response can be quickly sent.
	SetLibraryQueues([]LibraryQueue)

	// SetWaitingRunners stores an updated value that should be sent if a request to view
	// the waiting Runner is received.
	SetWaitingRunners(runnerNames []string)

	Start(ctx *context.Context, wg *sync.WaitGroup)
}

// The SettingsStorer defines how a struct which stores the settings in some manner
// should interact with other components of the application.
type SettingsStorer interface {
	Load() error
	Save() error
	Close() error

	// Getters and Setters

	HealthCheckInterval() uint64
	SetHealthCheckInterval(uint64)

	HealthCheckTimeout() uint64
	SetHealthCheckTimeout(uint64)

	LogVerbosity() string
	SetLogVerbosity(string)
}

type HealthCheckerDataStorer interface {
	DispatchedJobs() []DispatchedJob
	DeleteJob(uuid UUID) error
}

type LibraryManagerDataStorer interface {
	Libraries() []Library
	SaveLibrary(Library)

	FileEntryByPath(path string) File
	SaveFileEntry(f File)

	IsPathDispatched(path string) bool
}

type FileCacheDataStorer interface {
	Modtime(path string) (time.Time, error)
	Metadata(path string) (FileMetadata, error)

	SaveModtime(path string, t time.Time) error
	SaveMetadata(path string, f FileMetadata) error
}

type Logger interface {
	Trace(s string, i ...interface{})
	Debug(s string, i ...interface{})
	Info(s string, i ...interface{})
	Warn(s string, i ...interface{})
	Error(s string, i ...interface{})
	Critical(s string, i ...interface{})
}
