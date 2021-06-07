package controller

import (
	"context"
	"sync"
)

type mockHealthChecker struct {
	runCalled   bool
	startCalled bool
}

func (m *mockHealthChecker) Start(ctx *context.Context) {
	m.startCalled = true
}

func (m *mockHealthChecker) Run() (uuidsToNull []UUID) {
	m.runCalled = true
	return
}

type mockLibraryManager struct {
	importCalled            bool
	libSettingsCalled       bool
	libQueuesCalled         bool
	popJobCalled            bool
	updateLibSettingsCalled bool
	startCalled             bool
}

func (m *mockLibraryManager) Start(ctx *context.Context, wg *sync.WaitGroup) {
	m.startCalled = true
}

func (m *mockLibraryManager) ImportCompletedJobs([]Job) {
	m.importCalled = true
}

func (m *mockLibraryManager) LibrarySettings() (ls []Library) {
	m.libSettingsCalled = true
	return
}

func (m *mockLibraryManager) LibraryQueues() (lq []LibraryQueue) {
	m.libQueuesCalled = true
	return
}

func (m *mockLibraryManager) PopNewJob() (j Job) {
	m.popJobCalled = true
	return
}

func (m *mockLibraryManager) UpdateLibrarySettings(map[int]Library) {
	m.updateLibSettingsCalled = true
}

type mockRunnerCommunicator struct {
	completedJobsCalled  bool
	newJobCalled         bool
	needNewJobCalled     bool
	nullUUIDsCalled      bool
	waitingRunnersCalled bool
	startCalled          bool
}

func (m *mockRunnerCommunicator) Start(ctx *context.Context, wg *sync.WaitGroup) {
	m.startCalled = true
}

func (m *mockRunnerCommunicator) CompletedJobs() (j []Job) {
	m.completedJobsCalled = true
	return
}

func (m *mockRunnerCommunicator) NewJob(Job) {
	m.newJobCalled = true
}

func (m *mockRunnerCommunicator) NeedNewJob() bool {
	m.needNewJobCalled = true
	return true
}

func (m *mockRunnerCommunicator) NullifyUUIDs([]UUID) {
	m.nullUUIDsCalled = true
}

func (m *mockRunnerCommunicator) WaitingRunners() (runnerNames []string) {
	m.waitingRunnersCalled = true
	runnerNames = append(runnerNames, "TestRunner")
	return
}

type mockUserInterfacer struct {
	newLibSettingsCalled    bool
	setLibSettingsCalled    bool
	setLibQueuesCalled      bool
	setWaitingRunnersCalled bool
	startCalled             bool
}

func (m *mockUserInterfacer) Start(ctx *context.Context, wg *sync.WaitGroup) {
	m.startCalled = true
}

func (m *mockUserInterfacer) NewLibrarySettings() (ls map[int]Library) {
	m.newLibSettingsCalled = true
	return
}

func (m *mockUserInterfacer) SetLibrarySettings([]Library) {
	m.setLibSettingsCalled = true
}

func (m *mockUserInterfacer) SetLibraryQueues([]LibraryQueue) {
	m.setLibQueuesCalled = true
}

func (m *mockUserInterfacer) SetWaitingRunners(runnerNames []string) {
	m.setWaitingRunnersCalled = true
}

type mockLogger struct{}

func (m *mockLogger) Trace(s string, i ...interface{})    {}
func (m *mockLogger) Debug(s string, i ...interface{})    {}
func (m *mockLogger) Info(s string, i ...interface{})     {}
func (m *mockLogger) Warn(s string, i ...interface{})     {}
func (m *mockLogger) Error(s string, i ...interface{})    {}
func (m *mockLogger) Critical(s string, i ...interface{}) {}
