package controller

type UUID string

// The HealthChecker interface describes how a struct wishing to decide if a job's
// last update was long enough ago to mark the Runner doing it as unresponsive
// should interact with the Run function.
type HealthChecker interface {
	// Run loops through the provided slice of dispatched jobs and checks if any have
	// surpassed the allowed time between updates.
	Run() (uuidsToNull []UUID)
}

// The LibraryManager interface describes how a struct wishing to deal with user's
// libraries should interact with the Run function.
type LibraryManager interface {
	// LibraryQueues returns a slice of structs which represent the state of the queues
	LibraryQueues() (queues []struct{})
}

// The RunnerCommunicator interface describes how a struct wishing to communicate
// with external Runners should interact with the Run function.
type RunnerCommunicator interface {
	// NullifyUUIDs takes the provided slice of stringed UUIDs and marks them
	// so that if a Runner sends a request with a nullified UUID, it gets notified
	// that it is considered unresponsive and should acquire a new job.
	NullifyUUIDs(uuids []UUID)
}

// The UserInterfacer interface describes how a struct wishing to interact
// with the user should interact with the Run function.
type UserInterfacer interface {
}
