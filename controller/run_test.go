package controller

import (
	"context"
	"testing"
)

func TestRunFuncsCalled(t *testing.T) {
	ctx := context.Background()

	mockHealthChecker := MockHealthChecker{}
	mockLibraryManager := MockLibraryManager{}
	mockRunnerCommunicator := MockRunnerCommunicator{}
	mockUserInterfacer := MockUserInterfacer{}

	Run(&ctx, &mockHealthChecker, &mockLibraryManager, &mockRunnerCommunicator, &mockUserInterfacer, true)

	// Check that HealthChecker methods were run
	if !mockHealthChecker.runCalled {
		t.Errorf("HealthChecker.Run() wasn't called")
	}

	// Check that LibraryManager methods were run
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
