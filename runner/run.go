package runner

import (
	"context"
	"time"
)

// Run runs the basic loop of the Runner
func Run(ctx *context.Context, c Communicator, r CommandRunner) {
	for {
		if IsContextFinished(ctx) {
			break
		}
		// Send new job request
		ji, err := c.SendNewJobRequest(ctx)
		if err != nil {
			logger.Error(err.Error())
		}

		// Start job with request info
		r.Start(ji)

		// This allows us to rate limit to approx. 1 status POST request every 500ms, which keeps resources (tcp sockets) down.
		statusLastSent := time.Unix(0, 0)
		statusInterval := time.Duration(500 * time.Millisecond)
		sleepAmount := time.Duration(50 * time.Millisecond)

		for !r.Done() {
			// Rate limit how often we send status updates
			if time.Since(statusLastSent) < statusInterval {
				time.Sleep(sleepAmount)
				continue
			}
			statusLastSent = time.Now()

			// Get status from job
			status := r.Status()

			// Send status to Controller
			err = c.SendStatus(ctx, ji.UUID, status)
			if err != nil {
				logger.Error(err.Error())
			}

			if IsContextFinished(ctx) {
				break
			}
		}
		// If the context is finished, we want to avoid sending a misleading Job Complete request
		if IsContextFinished(ctx) {
			break
		}

		// Collect results from Command Runner
		cmdResults, failed, totalTime := r.Results()

		c.SendStatus(ctx, ji.UUID, JobStatus{
			Stage:                       "Copying to Controller",
			Percentage:                  "100",
			JobElapsedTime:              totalTime,
			FPS:                         "N/A",
			StageElapsedTime:            "N/A",
			StageEstimatedTimeRemaining: "N/A",
		})

		// Send job complete
		err = c.SendJobComplete(ctx, ji, failed, cmdResults)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func IsContextFinished(ctx *context.Context) bool {
	select {
	case <-(*ctx).Done():
		return true
	default:
		return false
	}
}
