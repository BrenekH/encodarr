package runner

import "context"

type mockCmdRunner struct {
	done           bool
	statusLoopout  int // How many calls to status should run before setting done to true
	jobInfo        JobInfo
	jobStatus      JobStatus
	commandResults CommandResults

	statusLoops int

	doneCalled    bool
	startCalled   bool
	statusCalled  bool
	resultsCalled bool
}

func (r *mockCmdRunner) Done() bool {
	r.doneCalled = true
	return r.done
}

func (r *mockCmdRunner) Start(ji JobInfo) {
	r.startCalled = true
	r.jobInfo = ji
}

func (r *mockCmdRunner) Status() JobStatus {
	r.statusCalled = true

	r.statusLoops++
	if r.statusLoops >= r.statusLoopout {
		r.done = true
	}
	return r.jobStatus
}

func (r *mockCmdRunner) Results() CommandResults {
	r.resultsCalled = true
	return r.commandResults
}

type mockCommunicator struct {
	jobReqJobInfo   JobInfo
	jobComJobInfo   JobInfo
	cmdResults      CommandResults
	statusUUID      string
	statusJobStatus JobStatus

	statusReturnErr error

	statusTimesCalled int

	jobCompleteCalled bool
	newJobCalled      bool
	statusCalled      bool

	// Used to to tell the mock to ignore any calls to SendStatus except for the very first one when setting statusUUID and statusJobStatus.
	statusSetOnlyFirst bool
}

func (c *mockCommunicator) SendJobComplete(ctx *context.Context, ji JobInfo, cr CommandResults) error {
	c.jobCompleteCalled = true
	c.jobComJobInfo = ji
	c.cmdResults = cr
	return nil
}

func (c *mockCommunicator) SendNewJobRequest(ctx *context.Context) (JobInfo, error) {
	c.newJobCalled = true
	return c.jobReqJobInfo, nil
}

func (c *mockCommunicator) SendStatus(ctx *context.Context, uuid string, js JobStatus) error {
	c.statusTimesCalled++
	if c.statusSetOnlyFirst && !c.statusCalled {
		c.statusUUID = uuid
		c.statusJobStatus = js
	}
	c.statusCalled = true
	return c.statusReturnErr
}
