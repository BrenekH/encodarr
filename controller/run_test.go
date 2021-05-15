package controller

import (
	"context"
	"testing"
)

func TestRunFuncsCalled(t *testing.T) {
	ctx := context.Background()

	mockLogger := MockLogger{}
	mockHealthChecker := MockHealthChecker{}
	mockLibraryManager := MockLibraryManager{}
	mockRunnerCommunicator := MockRunnerCommunicator{}
	mockUserInterfacer := MockUserInterfacer{}

	Run(&ctx, &mockLogger, &mockHealthChecker, &mockLibraryManager, &mockRunnerCommunicator, &mockUserInterfacer, true)

	// Check that HealthChecker methods were run
	if !mockHealthChecker.startCalled {
		t.Errorf("HealthChecker.Start() wasn't called")
	}
	if !mockHealthChecker.runCalled {
		t.Errorf("HealthChecker.Run() wasn't called")
	}

	// Check that LibraryManager methods were run
	if !mockLibraryManager.startCalled {
		t.Errorf("LibraryManager.Start() wasn't called")
	}
	if !mockLibraryManager.importCalled {
		t.Errorf("LibraryManager.ImportCompletedJobs() wasn't called")
	}
	if !mockLibraryManager.libSettingsCalled {
		t.Errorf("LibraryManager.LibrarySettings() wasn't called")
	}
	if !mockLibraryManager.libQueuesCalled {
		t.Errorf("LibraryManager.LibraryQueues() wasn't called")
	}
	if !mockLibraryManager.popJobCalled {
		t.Errorf("LibraryManager.PopNewJob() wasn't called")
	}
	if !mockLibraryManager.updateLibSettingsCalled {
		t.Errorf("LibraryManager.UpdateLibrarySettings wasn't called")
	}

	// Check that RunnerCommunicator methods were run
	if !mockRunnerCommunicator.startCalled {
		t.Errorf("RunnerCommunicator.Start() wasn't called")
	}
	if !mockRunnerCommunicator.completedJobsCalled {
		t.Errorf("RunnerCommunicator.CompletedJobs() wasn't called")
	}
	if !mockRunnerCommunicator.newJobCalled {
		t.Errorf("RunnerCommunicator.NewJob() wasn't called")
	}
	if !mockRunnerCommunicator.needNewJobCalled {
		t.Errorf("RunnerCommunicator.NeedNewJob() wasn't called")
	}
	if !mockRunnerCommunicator.nullUUIDsCalled {
		t.Errorf("RunnerCommunicator.NullifyUUIDs() wasn't called")
	}
	if !mockRunnerCommunicator.waitingRunnersCalled {
		t.Errorf("RunnerCommunicator.WaitingRunners() wasn't called")
	}

	// Check that UserInterfacer methods were run
	if !mockUserInterfacer.startCalled {
		t.Errorf("UserInterfacer.Start() wasn't called")
	}
	if !mockUserInterfacer.newLibSettingsCalled {
		t.Errorf("UserInterfacer.NewLibrarySettings() wasn't called")
	}
	if !mockUserInterfacer.setLibSettingsCalled {
		t.Errorf("UserInterfacer.SetLibrarySettings() wasn't called")
	}
	if !mockUserInterfacer.setLibQueuesCalled {
		t.Errorf("UserInterfacer.SetLibraryQueues() wasn't called")
	}
	if !mockUserInterfacer.setWaitingRunnersCalled {
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
