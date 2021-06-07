package runner_communicator

import (
	"context"
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

func NewRunnerHTTPApiV1(logger controller.Logger, httpServer controller.HTTPServer) RunnerHTTPApiV1 {
	return RunnerHTTPApiV1{logger: logger, httpServer: httpServer}
}

type RunnerHTTPApiV1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
}

func (r *RunnerHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	r.httpServer.Start(ctx, wg)

	// TODO: Add handlers to r.httpServer
}

func (r *RunnerHTTPApiV1) CompletedJobs() (j []controller.Job) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
	return
}

func (r *RunnerHTTPApiV1) NewJob(controller.Job) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
}

func (r *RunnerHTTPApiV1) NeedNewJob() (b bool) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
	return
}

func (r *RunnerHTTPApiV1) NullifyUUIDs([]controller.UUID) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
}

func (r *RunnerHTTPApiV1) WaitingRunners() (runnerNames []string) {
	r.logger.Critical("Not Implemented")
	// TODO: Implement
	return
}
