package runner

import (
	"context"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("Status Methods Called 5 Times", func(t *testing.T) {
		mCmdRunner := mockCmdRunner{
			done:          false,
			statusLoopout: 5,
		}
		mCommunicator := mockCommunicator{}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		if mCmdRunner.statusLoops != mCmdRunner.statusLoopout {
			t.Errorf("expected CmdRunner to run %v times but got %v instead", mCmdRunner.statusLoopout, mCmdRunner.statusLoops)
		}

		// We add one to the desired loops because the Run function sends a status to the Controller after the Command Runner is finished.
		if mCommunicator.statusTimesCalled != mCmdRunner.statusLoopout+1 {
			t.Errorf("expected Communicator to run %v times but got %v instead", mCmdRunner.statusLoopout+1, mCommunicator.statusTimesCalled)
		}
	})

	t.Run("Methods Called (Run Done Inner-loop)", func(t *testing.T) {
		mCmdRunner := mockCmdRunner{
			done:          false,
			statusLoopout: 1,
		}
		mCommunicator := mockCommunicator{}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		// Command Runner
		if !mCmdRunner.doneCalled {
			t.Errorf("expected Done to be called, but it wasn't")
		}

		if !mCmdRunner.statusCalled {
			t.Errorf("expected Status to be called, but it wasn't")
		}

		if !mCmdRunner.startCalled {
			t.Errorf("expected Start to be called, but it wasn't")
		}

		if !mCmdRunner.resultsCalled {
			t.Errorf("expected Results to be called, but it wasn't")
		}

		// Communicator
		if !mCommunicator.jobCompleteCalled {
			t.Errorf("expected SendJobComplete to be called, but it wasn't")
		}

		if !mCommunicator.newJobCalled {
			t.Errorf("expected SendNewJobRequest to be called, but it wasn't")
		}

		if !mCommunicator.statusCalled {
			t.Errorf("expected SendStatus to be called, but it wasn't")
		}
	})

	t.Run("Methods Called (Don't Run Done Inner-loop)", func(t *testing.T) {
		mCmdRunner := mockCmdRunner{
			done: true,
		}
		mCommunicator := mockCommunicator{}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		// Command Runner
		if !mCmdRunner.doneCalled {
			t.Errorf("expected Done to be called, but it wasn't")
		}

		if mCmdRunner.statusCalled {
			t.Errorf("expected Status to not be called, but it was")
		}

		if !mCmdRunner.startCalled {
			t.Errorf("expected Start to be called, but it wasn't")
		}

		if !mCmdRunner.resultsCalled {
			t.Errorf("expected Results to be called, but it wasn't")
		}

		// Communicator
		if !mCommunicator.jobCompleteCalled {
			t.Errorf("expected SendJobComplete to be called, but it wasn't")
		}

		if !mCommunicator.newJobCalled {
			t.Errorf("expected SendNewJobRequest to be called, but it wasn't")
		}

		if !mCommunicator.statusCalled {
			t.Errorf("expected SendStatus to be called, but it wasn't")
		}
	})

	t.Run("JobInfo from SendNewJobRequest is Received Properly by Start", func(t *testing.T) {
		ji := JobInfo{
			UUID:        "uuid-4",
			File:        "/library/1/my_file.mkv",
			InFile:      "/tmp/in.mkv",
			OutFile:     "/tmp/out.mkv",
			CommandArgs: []string{"-v:c", "hevc"},
			MediaInfo:   MediaInfo{},
		}
		mCmdRunner := mockCmdRunner{
			done:          false,
			statusLoopout: 2,
		}
		mCommunicator := mockCommunicator{
			jobReqJobInfo: ji,
		}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		if !reflect.DeepEqual(ji, mCmdRunner.jobInfo) {
			t.Errorf("expected %v but got %v", ji, mCmdRunner.jobInfo)
		}
	})

	t.Run("JobStatus from Status is Received Properly by SendStatus", func(t *testing.T) {
		js := JobStatus{
			Stage:                       "Testing",
			Percentage:                  "32.0",
			JobElapsedTime:              "1s",
			FPS:                         "9",
			StageElapsedTime:            "1s",
			StageEstimatedTimeRemaining: "9s",
		}
		mCmdRunner := mockCmdRunner{
			done:          false,
			statusLoopout: 2,
			jobStatus:     js,
		}
		mCommunicator := mockCommunicator{
			statusSetOnlyFirst: true,
		}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		if !reflect.DeepEqual(js, mCommunicator.statusJobStatus) {
			t.Errorf("expected %v but got %v", js, mCommunicator.statusJobStatus)
		}
	})

	t.Run("CommandResults from Results are Received Properly by SendJobComplete", func(t *testing.T) {
		cr := CommandResults{}
		mCmdRunner := mockCmdRunner{
			done:           false,
			statusLoopout:  2,
			commandResults: cr,
		}
		mCommunicator := mockCommunicator{}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		if !reflect.DeepEqual(cr, mCommunicator.cmdResults) {
			t.Errorf("expected %v but got %v", cr, mCommunicator.cmdResults)
		}
	})

	t.Run("Unresponsive Err from SendStatus Doesn't Send Job Complete", func(t *testing.T) {
		mCmdRunner := mockCmdRunner{
			done:          false,
			statusLoopout: 2,
		}
		mCommunicator := mockCommunicator{
			statusReturnErr: ErrUnresponsive,
		}
		ctx := context.Background()

		Run(&ctx, &mCommunicator, &mCmdRunner, true)

		if mCommunicator.jobCompleteCalled {
			t.Errorf("SendJobComplete was unexpectedly called")
		}
	})
}

func TestIsContextFinished(t *testing.T) {
	fCtx, cancel := context.WithCancel(context.Background())
	cancel()
	uCtx := context.Background()

	tests := []struct {
		name string
		in   context.Context
		out  bool
	}{
		{name: "Finished Context", in: fCtx, out: true},
		{name: "Unfinished Context", in: uCtx, out: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := IsContextFinished(&test.in)

			if out != test.out {
				t.Errorf("expected %v but got %v", test.out, out)
			}
		})
	}
}
