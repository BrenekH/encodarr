package controller

import "time"

type UUID string

// Job represents a job to be carried out by a Runner.
type Job struct {
	UUID     UUID     `json:"uuid"`
	Path     string   `json:"path"`
	Command  []string `json:"command"`
	Metadata struct{} `json:"metadata"` // TODO: Define metadata as its own struct
}

// Library represents a single library.
type Library struct {
	ID              int           `json:"id"`
	Folder          string        `json:"folder"`
	Priority        int           `json:"priority"`
	FsCheckInterval time.Duration `json:"fs_check_interval"`
	Queue           LibraryQueue  `json:"queue"`
	PathMasks       []string      `json:"path_masks"`

	// TODO: Figure out how to replace PluginPipeline
	// Pipeline        PluginPipeline `json:"pipeline"`
}

// LibraryQueue represents a singular queue belonging to one library.
type LibraryQueue struct{}

type DispatchedJob struct {
	UUID        UUID      `json:"uuid"`
	Runner      string    `json:"runner"`
	Job         Job       `json:"job"`
	Status      JobStatus `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}
