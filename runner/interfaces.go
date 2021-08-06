package runner

import "context"

// Communicator defines how a struct which talks with a Controller should behave.
type Communicator interface {
	SendJobComplete(*context.Context, JobInfo, CommandResults) error
	SendNewJobRequest(*context.Context) (JobInfo, error)
	SendStatus(*context.Context, string, JobStatus) error
}

// CommandRunner defines how a struct which runs the FFmpeg commands should behave.
type CommandRunner interface {
	Done() bool
	Start(JobInfo)
	Status() JobStatus
	Results() CommandResults
}
