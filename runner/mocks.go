package runner

type MockCmdRunner struct{}

func (r *MockCmdRunner) Done() bool {
	return true
}

func (r *MockCmdRunner) Start(s string) {}

func (r *MockCmdRunner) Status() {}
