package controller

import (
	"context"
	"sync"
)

// Run is the "top-level" function for running the Encodarr Controller. It calls all of the injected
// dependencies in order to operate.
func Run(ctx *context.Context, logger Logger, hc HealthChecker, lm LibraryManager, rc RunnerCommunicator, ui UserInterfacer, testMode bool) {
	wg := sync.WaitGroup{}
	hc.Start(ctx)
	lm.Start(ctx, &wg)
	rc.Start(ctx, &wg)
	ui.Start(ctx, &wg)
	looped := false

	for {
		// A while loop will skip if its condition is false even on the first run.
		// Using the looped var allows a do-while run for testing.
		if testMode && looped {
			break
		}
		if testMode {
			looped = true
		}

		if IsContextFinished(ctx) {
			break
		}

		// Run health check and null any unresponsive Runners
		uuidsToNull := hc.Run()
		rc.NullifyUUIDs(uuidsToNull)

		// Update the UserInterfacer library settings cache
		ls := lm.LibrarySettings()
		ui.SetLibrarySettings(ls)

		// Apply user changes to library settings
		lsUserChanges := ui.NewLibrarySettings()
		lm.UpdateLibrarySettings(lsUserChanges)

		// TODO: Remove the setting of queues in the UserInterfacer
		// Update queues cache in the UserInterfacer
		lq := lm.LibraryQueues()
		ui.SetLibraryQueues(lq)

		// Update waiting runners to be shown to the user
		wr := rc.WaitingRunners()
		ui.SetWaitingRunners(wr)

		// Send new job to the RunnerCommunicator if there is a waiting Runner
		if rc.NeedNewJob() {
			nj := lm.PopNewJob()
			rc.NewJob(nj)
		}

		// Import completed jobs
		cj := rc.CompletedJobs()
		lm.ImportCompletedJobs(cj)
	}

	// Wait for goroutines to shut down
	wg.Wait()
}
