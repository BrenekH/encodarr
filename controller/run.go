package controller

import "context"

// Run is the "top-level" function for running the Encodarr Controller. It calls all of the injected
// dependencies in order to operate.
func Run(ctx *context.Context, ds DataStorer, hc HealthChecker, lm LibraryManager, rc RunnerCommunicator, ui UserInterfacer) {
	for {
		if IsContextFinished(ctx) {
			break
		}

		// GetUserInput returns any settings that the user has changed. Other input may be available in the future, so plan for that.
		ui.GetUserInput()

		// Launches a goroutine for libraries that are due for a file system check.
		lm.StartFSChecks(ctx)

		// Loops through the dispatched jobs, and determines if any are unresponsive.
		hc.RunCheck([]struct{}{})

		// Remove any dispatched jobs that have become unresponsive.
		ds.RemoveDispatchedJob([]string{})

		// Adds UUIDs to null from the health check.
		rc.AddNullUUIDs([]string{})

		// Returns Runner statuses to be saved and communicated to the user.
		rc.GetRunnerStatuses()

		// Takes the statuses grabbed from GetRunnerInfo and allows them to be seen by the user.
		ui.SetJobStatuses([]struct{}{})

		// Takes the new job requests from GetRunnerInfo and sends them to the user for display as waiting Runners.
		ui.SetWaitingRunners([]string{})

		//? What to do with completed jobs?

		//? Should the CommandDecider (plugins and naive approach(current way)) be a top-level interface?
		//?   - Passing data from the file system checks up to the Run func, to be passed into CommandDecider.
	}
}
