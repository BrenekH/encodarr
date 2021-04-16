package mock

type CmdRunner struct{}

func (r *CmdRunner) Done() bool {
	return true
}

func (r *CmdRunner) Start(s string) {}

func (r *CmdRunner) Status() {}
