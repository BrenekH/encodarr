package controller

import "context"

func Run(ctx *context.Context, ds DataStorer, hc HealthChecker, lm LibraryManager, rc RunnerCommunicator, ui UserInterfacer) {
	for {
		if IsContextFinished(ctx) {
			break
		}

		// GetUserInput returns any settings that the user has changed. Other input may be available in the future, so plan for that.
		ui.GetUserInput()

		// Launches a goroutine for library that is due for a file system check.
		lm.StartFSChecks(ctx)

		// Loops through the dispatched jobs, and determines if any are unresponsive.
		hc.RunCheck()

		// Remove any dispatched jobs that have become unresponsive.
		ds.RemoveDispatchedJob([]string{})

		// Adds UUIDs to null from the health check.
		rc.AddNullUUIDs([]string{})

		// Returns Runner statuses to be saved and communicated to the user.
		rc.GetRunnerInfo()

		// Takes the statuses grabbed from GetRunnerInfo and allows them to be seen by the user.
		ui.SetRunnerStatuses([]struct{}{})

		// Takes the new job requests from GetRunnerInfo and sends them to the user for display as waiting Runners.
		ui.SetWaitingRunners([]string{})
	}
}
