package sqlite

import (
	"encoding/json"

	"github.com/BrenekH/encodarr/controller"
)

func NewRunnerCommunicatorAdapter(db *SQLiteDatabase, logger controller.Logger) RunnerCommunicatorAdapater {
	return RunnerCommunicatorAdapater{db: db, logger: logger}
}

type RunnerCommunicatorAdapater struct {
	db     *SQLiteDatabase
	logger controller.Logger
}

func (r *RunnerCommunicatorAdapater) DispatchedJob(uuid controller.UUID) (controller.DispatchedJob, error) {
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

func (r *RunnerCommunicatorAdapater) SaveDispatchedJob(dJob controller.DispatchedJob) error {
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
