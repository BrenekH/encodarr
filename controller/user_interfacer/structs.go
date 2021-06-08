package user_interfacer

import "github.com/BrenekH/encodarr/controller"

type runningJSONResponse struct {
	DispatchedJobs []filteredDispatchedJob `json:"jobs"`
}

type filteredDispatchedJob struct {
	Job        filteredJob          `json:"job"`
	RunnerName string               `json:"runner_name"`
	Status     controller.JobStatus `json:"status"`
}

type filteredJob struct {
	UUID    controller.UUID `json:"uuid"`
	Path    string          `json:"path"`
	Command []string        `json:"command"`
}

type settingsJSON struct {
	FileSystemCheckInterval string
	HealthCheckInterval     string
	HealthCheckTimeout      string
	LogVerbosity            string
}

type humanizedHistoryEntry struct {
	File              string   `json:"file"`
	DateTimeCompleted string   `json:"datetime_completed"`
	Warnings          []string `json:"warnings"`
	Errors            []string `json:"errors"`
}

type historyJSON struct {
	History []humanizedHistoryEntry `json:"history"`
}
