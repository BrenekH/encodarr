package controller

import (
	"context"
	"sync"
	"time"
)

// Run is the "top-level" function for running the Encodarr Controller. It calls all of the injected
// dependencies in order to operate.
func Run(ctx *context.Context, logger Logger, hc HealthChecker, lm LibraryManager, rc RunnerCommunicator, ui UserInterfacer, setLogLvl func(), testMode bool) {
	wg := sync.WaitGroup{}
	hc.Start(ctx)
	lm.Start(ctx, &wg)
	rc.Start(ctx, &wg)
	ui.Start(ctx, &wg)
	looped := false

	loopsPerSec := 20
	ticker := time.NewTicker(time.Second / time.Duration(loopsPerSec))

	for range ticker.C {
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
		if ls, err := lm.LibrarySettings(); err == nil {
			ui.SetLibrarySettings(ls)
		}

		// Apply user changes to library settings
		lsUserChanges := ui.NewLibrarySettings()
		lm.UpdateLibrarySettings(lsUserChanges)

		// Update waiting runners to be shown to the user
		wr := rc.WaitingRunners()
		ui.SetWaitingRunners(wr)

		// Send new job to the RunnerCommunicator if there is a waiting Runner
		if rc.NeedNewJob() {
			if nj, err := lm.PopNewJob(); err == nil {
				rc.NewJob(nj)
			}
		}

		// Import completed jobs
		cj := rc.CompletedJobs()
		lm.ImportCompletedJobs(cj)

		// Apply the log level to the actual handler
		setLogLvl()
	}

	// Wait for goroutines to shut down
	wg.Wait()
}
