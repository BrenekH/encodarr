package controller

import "context"

// The DataStorer interface describes how a struct wishing to store information
// such as configuration settings and library/dispatched job data should interact
// with the Run function.
type DataStorer interface {
	// RemoveDispatchedJob removes any dispatched jobs that match the passed UUIDs.
	RemoveDispatchedJob(UUIDs []string)
}

// The HealthChecker interface describes how a struct wishing to decide if a job's
// last update was long enough ago to mark the Runner doing it as unresponsive
// should interact with the Run function.
type HealthChecker interface {
	// RunCheck loops through the provided slice of dispatched jobs and checks if any have
	// surpassed the allowed time between updates.
	RunCheck([]struct{}) []string
}

// The LibraryManager interface describes how a struct wishing to deal with user's
// libraries should interact with the Run function.
type LibraryManager interface {
	// StartFSChecks identifies which libraries are due for an update and then
	// spawns a goroutine for all that do require one.
	StartFSChecks(ctx *context.Context)
}

// The RunnerCommunicator interface describes how a struct wishing to communicate
// with external Runners should interact with the Run function.
type RunnerCommunicator interface {
	// AddNullUUIDs takes the provided slice of stringed UUIDs and marks them
	// so that if a Runner sends a request with a nullified UUID, it gets notified
	// that it is considered unresponsive and should acquire a new job.
	AddNullUUIDs(uuids []string)

	// GetRunnerStatuses returns the updated statuses of all active Runners.
	GetRunnerStatuses()
}

// The UserInterfacer interface describes how a struct wishing to interact
// with the user should interact with the Run function.
type UserInterfacer interface {
	// GetUserInput returns a struct with information about what the user has requested
	// since the last time this function was called.
	GetUserInput()

	// SetJobStatuses takes in a slice with information about each active job's status
	// so that it can be displayed to the user.
	SetJobStatuses(statuses []struct{})

	// SetWaitingRunners takes in a slice which directly corresponds to Runners' names
	// which have requested a job, but not yet received one.
	SetWaitingRunners(runnerNames []string)
}
