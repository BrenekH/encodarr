package user_interfacer

import "github.com/BrenekH/encodarr/controller"

func filterDispatchedJobs(dJobs []controller.DispatchedJob) []filteredDispatchedJob {
	fDJobs := make([]filteredDispatchedJob, 0)
	for _, dJob := range dJobs {
		fDJobs = append(fDJobs, filteredDispatchedJob{
			Job: filteredJob{
				UUID:    dJob.Job.UUID,
				Path:    dJob.Job.Path,
				Command: dJob.Job.Command,
			},
			RunnerName: dJob.Runner,
			Status:     dJob.Status,
		})
	}
	return fDJobs
}
