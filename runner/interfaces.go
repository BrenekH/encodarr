package runner

import "context"

type Communicator interface {
	SendJobComplete(*context.Context, JobInfo, CommandResults) error
	SendNewJobRequest(*context.Context) (JobInfo, error)
	SendStatus(*context.Context, string, JobStatus) error
}

type CommandRunner interface {
	Done() bool
	Start(JobInfo)
	Status() JobStatus
	Results() CommandResults
}
