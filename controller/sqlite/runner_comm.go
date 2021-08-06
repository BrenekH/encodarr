package sqlite

import (
	"encoding/json"

	"github.com/BrenekH/encodarr/controller"
)

// NewRunnerCommunicatorAdapter returns an instantiated RunnerCommunicatorAdapter.
func NewRunnerCommunicatorAdapter(db *Database, logger controller.Logger) RunnerCommunicatorAdapter {
	return RunnerCommunicatorAdapter{db: db, logger: logger}
}

// RunnerCommunicatorAdapter is a struct that satisfies the interface that connects a RunnerCommunicator
// to a storage medium.
type RunnerCommunicatorAdapter struct {
	db     *Database
	logger controller.Logger
}

// DispatchedJob uses the provided uuid to retrieve a dispatched job from the database.
func (r *RunnerCommunicatorAdapter) DispatchedJob(uuid controller.UUID) (controller.DispatchedJob, error) {
	row := r.db.Client.QueryRow("SELECT job, status, runner, last_updated FROM dispatched_jobs WHERE uuid = $1;", uuid)

	d := controller.DispatchedJob{UUID: uuid}
	bJob := []byte{}
	bStatus := []byte{}

	row.Scan(
		&bJob,
		&bStatus,
		&d.Runner,
		&d.LastUpdated,
	)

	if err := json.Unmarshal(bJob, &d.Job); err != nil {
		return d, err
	}

	if err := json.Unmarshal(bStatus, &d.Status); err != nil {
		return d, err
	}

	return d, nil
}

// SaveDispatchedJob saves the provided dispatched job to the database.
func (r *RunnerCommunicatorAdapter) SaveDispatchedJob(dJob controller.DispatchedJob) error {
	bJob, err := json.Marshal(dJob.Job)
	if err != nil {
		return err
	}

	bStatus, err := json.Marshal(dJob.Status)
	if err != nil {
		return err
	}

	_, err = r.db.Client.Exec("INSERT INTO dispatched_jobs (uuid, job, status, runner, last_updated) VALUES ($1, $2, $3, $4, $5) ON CONFLICT(uuid) DO UPDATE SET uuid=$1, job=$2, status=$3, runner=$4, last_updated=$5;",
		dJob.UUID,
		bJob,
		bStatus,
		dJob.Runner,
		dJob.LastUpdated,
	)
	return err
}
