package controller

import "context"

type DataStorer interface {
	RemoveDispatchedJob(UUIDs []string)
}

type HealthChecker interface {
	RunCheck()
}

type LibraryManager interface {
	StartFSChecks(ctx *context.Context)
}

type RunnerCommunicator interface {
	AddNullUUIDs([]string)
	GetRunnerInfo()
}

type UserInterfacer interface {
	GetUserInput()
	SetRunnerStatuses(statuses []struct{})
	SetWaitingRunners(runnerNames []string)
}
