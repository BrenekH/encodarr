package runner

import "time"

type JobInfo struct {
	UUID          string
	File          string
	InFile        string
	OutFile       string
	CommandArgs   []string
	MediaDuration float32
}

type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}

type CommandResults struct {
	Failed         bool
	JobElapsedTime time.Duration
	Warnings       []string
	Errors         []string
}
