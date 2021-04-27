package controller

import "context"

func Run(ctx *context.Context) {
	for {
		if IsContextFinished(ctx) {
			break
		}
		// UserInterfacer.GetUserInput()
		//   GetUserInput returns any settings that the user has changed. Other input may be available in the future, so plan for that.

		// StartLibraryChecks(ctx)
		//   Launches a goroutine for library that is due for a file system check.

		// RunHealthCheck()
		//   Loops through the dispatched jobs, and determines if any are unresponsive.

		// DataStorer.RemoveDispatchedJob([]uuid)
		//   Remove any dispatched jobs that have become unresponsive.

		// RunnerCommunicator.GetRunnerInfo()
		//   Returns Runner statuses to be saved and communicated to the user.

		// UserInterfacer.SetRunnerStatuses([]UnnamedStruct)
		//   Takes the statuses grabbed from GetRunnerInfo and allows them to be seen by the user.

		// UserInterfacer.SetWaitingRunners([]string)
		//   Takes the new job requests from GetRunnerInfo and sends them to the user for display as waiting Runners.
	}
}
