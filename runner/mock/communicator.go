package mock

import (
	"context"

	"github.com/BrenekH/encodarr/runner"
)

type MockCommunicator struct{}

func (c *MockCommunicator) SendJobComplete(ctx *context.Context, ji runner.JobInfo, cr runner.CommandResults) error {
	return nil
}

func (c *MockCommunicator) SendNewJobRequest(ctx *context.Context) (runner.JobInfo, error) {
	return runner.JobInfo{}, nil
}

func (c *MockCommunicator) SendStatus(ctx *context.Context, uuid string, js runner.JobStatus) error {
	return nil
}
