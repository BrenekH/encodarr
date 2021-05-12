package controller

type MockHealthChecker struct {
	runCalled bool
}

func (m *MockHealthChecker) Run() (uuidsToNull []UUID) {
	m.runCalled = true
	return
}

type MockLibraryManager struct {
	importCalled            bool
	libSettingsCalled       bool
	libQueuesCalled         bool
	popJobCalled            bool
	updateLibSettingsCalled bool
}

func (m *MockLibraryManager) ImportCompletedJobs([]Job) {
	m.importCalled = true
}

func (m *MockLibraryManager) LibrarySettings() (ls []LibrarySettings) {
	m.libSettingsCalled = true
	return
}

func (m *MockLibraryManager) LibraryQueues() (lq []LibraryQueue) {
	m.libQueuesCalled = true
	return
}

func (m *MockLibraryManager) PopNewJob() (j Job) {
	m.popJobCalled = true
	return
}

func (m *MockLibraryManager) UpdateLibrarySettings(map[string]LibrarySettings) {
	m.updateLibSettingsCalled = true
}

type MockRunnerCommunicator struct {
	completedJobsCalled  bool
	newJobCalled         bool
	nullUUIDsCalled      bool
	waitingRunnersCalled bool
}

func (m *MockRunnerCommunicator) CompletedJobs() (j []Job) {
	m.completedJobsCalled = true
	return
}

func (m *MockRunnerCommunicator) NewJob(Job) {
	m.newJobCalled = true
}

func (m *MockRunnerCommunicator) NullifyUUIDs([]UUID) {
	m.nullUUIDsCalled = true
}

func (m *MockRunnerCommunicator) WaitingRunners() (runnerNames []string) {
	m.waitingRunnersCalled = true
	runnerNames = append(runnerNames, "TestRunner")
	return
}

type MockUserInterfacer struct {
	newLibSettingsCalled    bool
	setLibSettingsCalled    bool
	setLibQueuesCalled      bool
	setWaitingRunnersCalled bool
}

func (m *MockUserInterfacer) NewLibrarySettings() (ls map[string]LibrarySettings) {
	m.newLibSettingsCalled = true
	return
}

func (m *MockUserInterfacer) SetLibrarySettings([]LibrarySettings) {
	m.setLibSettingsCalled = true
}

func (m *MockUserInterfacer) SetLibraryQueues([]LibraryQueue) {
	m.setLibQueuesCalled = true
}

func (m *MockUserInterfacer) SetWaitingRunners(runnerNames []string) {
	m.setWaitingRunnersCalled = true
}
