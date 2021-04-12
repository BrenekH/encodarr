package runner

import "context"

type Communicator interface {
	SendJobComplete(*context.Context) error
	SendNewJobRequest(*context.Context) (JobInfo, error)
	SendStatus(*context.Context) error
}

type CommandRunner interface {
	Done() bool
	Start([]string)
	Status()
}

type JobInfo struct {
	CommandArgs []string
}
