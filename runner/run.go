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
			continue
		}

		// Start job with request info
		r.Start(ji)

		// This allows us to rate limit to approx. 1 status POST request every 500ms, which keeps resources (tcp sockets) down.
		statusLastSent := time.Unix(0, 0)
		statusInterval := time.Duration(500 * time.Millisecond)
		sleepAmount := time.Duration(50 * time.Millisecond)

		unresponsive := false

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
				if err == ErrUnresponsive {
					logger.Warn(err.Error())
					unresponsive = true
					break
				} else {
					logger.Error(err.Error())
				}
			}

			if IsContextFinished(ctx) {
				break
			}
		}
		// If we are detected as unresponsive, skip sending the job complete request.
		if unresponsive {
			continue
		}

		// If the context is finished, we want to avoid sending a misleading Job Complete request
		if IsContextFinished(ctx) {
			break
		}

		// Collect results from Command Runner
		cmdResults := r.Results()

		// Make sure that the Web UI properly states that we are copying the result to the Controller.
		// Setting Percentage to 100 also makes sure that the Runner card appears at the top of the page.
		err = c.SendStatus(ctx, ji.UUID, JobStatus{
			Stage:                       "Copying to Controller",
			Percentage:                  "100",
			JobElapsedTime:              cmdResults.JobElapsedTime.String(),
			FPS:                         "N/A",
			StageElapsedTime:            "N/A",
			StageEstimatedTimeRemaining: "N/A",
		})
		if err != nil {
			logger.Warn(err.Error())
		}

		// Send job complete
		err = c.SendJobComplete(ctx, ji, cmdResults)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

// IsContextFinshed returns a boolean indicating whether or not a context.Context is finished.
// This replaces the need to use a select code block.
func IsContextFinished(ctx *context.Context) bool {
	select {
	case <-(*ctx).Done():
		return true
	default:
		return false
	}
}
