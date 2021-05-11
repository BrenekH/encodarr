package controller

import "context"

// Run is the "top-level" function for running the Encodarr Controller. It calls all of the injected
// dependencies in order to operate.
func Run(ctx *context.Context, hc HealthChecker, lm LibraryManager, rc RunnerCommunicator, ui UserInterfacer) {

}
