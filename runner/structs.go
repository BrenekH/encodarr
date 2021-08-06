package runner

import "time"

// JobInfo defines the information about a Job that is sent from the Controller.
type JobInfo struct {
	UUID          string
	File          string
	InFile        string
	OutFile       string
	CommandArgs   []string
	MediaDuration float32
}

// JobStatus defines the information to be reported about the current state of a running job.
type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}

// CommandResults is the results of the FFmpeg command that the Runner was told by the Controller to run.
type CommandResults struct {
	Failed         bool
	JobElapsedTime time.Duration
	Warnings       []string
	Errors         []string
}
