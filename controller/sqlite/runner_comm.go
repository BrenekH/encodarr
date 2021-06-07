package sqlite

import "github.com/BrenekH/encodarr/controller"

func NewRunnerCommunicatorAdapter(db *SQLiteDatabase, logger controller.Logger) RunnerCommunicatorAdapater {
	return RunnerCommunicatorAdapater{db: db, logger: logger}
}

type RunnerCommunicatorAdapater struct {
	db     *SQLiteDatabase
	logger controller.Logger
}

func (r *RunnerCommunicatorAdapater) DispatchedJob() (controller.DispatchedJob, error) {
	r.logger.Critical("Not yet implemented")
	// TODO: Implement
	return controller.DispatchedJob{}, nil
}

func (r *RunnerCommunicatorAdapater) SaveDispatchedJob(dJob controller.DispatchedJob) error {
	r.logger.Critical("Not yet implemented")
	// TODO: Implement
	return nil
}
