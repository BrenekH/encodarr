package job_health

import (
	"context"

	"github.com/BrenekH/encodarr/controller"
)

func NewChecker(ds controller.HealthCheckerDataStorer) Checker {
	return Checker{ds: &ds}
}

type Checker struct {
	ds *controller.HealthCheckerDataStorer
}

// Run loops through the provided slice of dispatched jobs and checks if any have
// surpassed the allowed time between updates.
func (c *Checker) Run() (uuidsToNull []controller.UUID) {
	// TODO: Implement
	return
}

// Start just satisfies the controller.HealthChecker interface.
// There is no implemented functionality.
func (c *Checker) Start(ctx *context.Context) {}
