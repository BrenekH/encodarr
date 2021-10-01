package controller

import (
	"context"
	"testing"
)

func TestRunFuncsCalled(t *testing.T) {
	ctx := context.Background()

	mLogger := mockLogger{}
	mHealthChecker := mockHealthChecker{}
	mLibraryManager := mockLibraryManager{}
	mRunnerCommunicator := mockRunnerCommunicator{}
	mUserInterfacer := mockUserInterfacer{}

	Run(&ctx, &mLogger, &mHealthChecker, &mLibraryManager, &mRunnerCommunicator, &mUserInterfacer, func() {}, true)

	// Check that HealthChecker methods were run
	if !mHealthChecker.startCalled {
		t.Errorf("HealthChecker.Start() wasn't called")
	}
	if !mHealthChecker.runCalled {
		t.Errorf("HealthChecker.Run() wasn't called")
	}

	// Check that LibraryManager methods were run
	if !mLibraryManager.startCalled {
		t.Errorf("LibraryManager.Start() wasn't called")
	}
	if !mLibraryManager.importCalled {
		t.Errorf("LibraryManager.ImportCompletedJobs() wasn't called")
	}
	if !mLibraryManager.libSettingsCalled {
		t.Errorf("LibraryManager.LibrarySettings() wasn't called")
	}
	if !mLibraryManager.popJobCalled {
		t.Errorf("LibraryManager.PopNewJob() wasn't called")
	}
	if !mLibraryManager.updateLibSettingsCalled {
		t.Errorf("LibraryManager.UpdateLibrarySettings wasn't called")
	}

	// Check that RunnerCommunicator methods were run
	if !mRunnerCommunicator.startCalled {
		t.Errorf("RunnerCommunicator.Start() wasn't called")
	}
	if !mRunnerCommunicator.completedJobsCalled {
		t.Errorf("RunnerCommunicator.CompletedJobs() wasn't called")
	}
	if !mRunnerCommunicator.newJobCalled {
		t.Errorf("RunnerCommunicator.NewJob() wasn't called")
	}
	if !mRunnerCommunicator.needNewJobCalled {
		t.Errorf("RunnerCommunicator.NeedNewJob() wasn't called")
	}
	if !mRunnerCommunicator.nullUUIDsCalled {
		t.Errorf("RunnerCommunicator.NullifyUUIDs() wasn't called")
	}
	if !mRunnerCommunicator.waitingRunnersCalled {
		t.Errorf("RunnerCommunicator.WaitingRunners() wasn't called")
	}

	// Check that UserInterfacer methods were run
	if !mUserInterfacer.startCalled {
		t.Errorf("UserInterfacer.Start() wasn't called")
	}
	if !mUserInterfacer.newLibSettingsCalled {
		t.Errorf("UserInterfacer.NewLibrarySettings() wasn't called")
	}
	if !mUserInterfacer.setLibSettingsCalled {
		t.Errorf("UserInterfacer.SetLibrarySettings() wasn't called")
	}
	if !mUserInterfacer.setWaitingRunnersCalled {
		t.Errorf("UserInterfacer.SetWaitingRunners() wasn't called")
	}
}

// Test to write
//   - rc.NullifyUUIDs() is called with the return value of hc.Run()
//   - ui.SetLibrarySettings() is called with the return value of lm.LibrarySettings()
//   - lm.UpdateLibrarySettings() is called with the return value of ui.NewLibrarySettings()
//   - ui.SetLibraryQueues() is called with the return value of lm.LibraryQueues()
//   - ui.SetWaitingRunners() is called with the return value of rc.WaitingRunners()
//   - rc.NewJob() is called with the return value of lm.PopNewJob() only when rc.NeedNewJob() returns true
//   - lm.ImportCompletedJobs() is called with the return value of rc.CompletedJobs()
//   - run breaks when context is done (test mode is false and context canceled after a few milliseconds)
