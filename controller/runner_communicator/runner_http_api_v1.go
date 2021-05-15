package runner_communicator

import (
	"context"
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

func NewRunnerHTTPApiV1(logger controller.Logger) RunnerHTTPApiV1 {
	return RunnerHTTPApiV1{logger: logger}
}

type RunnerHTTPApiV1 struct {
	logger controller.Logger
}

func (r *RunnerHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	r.logger.Critical("Not Implemented")
	// Run r.httpServer.Start(ctx, wg)
	// Add handlers to r.httpServer
}

func (r *RunnerHTTPApiV1) CompletedJobs() (j []controller.Job) {
	r.logger.Critical("Not Implemented")
	return
}

func (r *RunnerHTTPApiV1) NewJob(controller.Job) {
	r.logger.Critical("Not Implemented")
}

func (r *RunnerHTTPApiV1) NeedNewJob() (b bool) {
	r.logger.Critical("Not Implemented")
	return
}

func (r *RunnerHTTPApiV1) NullifyUUIDs([]controller.UUID) {
	r.logger.Critical("Not Implemented")
}

func (r *RunnerHTTPApiV1) WaitingRunners() (runnerNames []string) {
	r.logger.Critical("Not Implemented")
	return
}
