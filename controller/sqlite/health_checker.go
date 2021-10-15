package sqlite

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

// NewHealthCheckerAdapter returns a new instantiated HealthCheckerAdapter.
func NewHealthCheckerAdapter(db *Database, logger controller.Logger) HealthCheckerAdapter {
	return HealthCheckerAdapter{db: db, logger: logger}
}

// HealthCheckerAdapter satisfies the controller.HealthCheckerDataStorer interface by turning interface
// requests into SQL requests that are passed on to an underlying SQLiteDatabase.
type HealthCheckerAdapter struct {
	db     *Database
	logger controller.Logger
}

// DispatchedJobs returns all of the dispatched jobs in the database.
func (h *HealthCheckerAdapter) DispatchedJobs() []controller.DispatchedJob {
	returnSlice := make([]controller.DispatchedJob, 0)

	rows, err := h.db.Client.Query("SELECT uuid, runner, job, status, last_updated FROM dispatched_jobs;")
	if err != nil {
		h.logger.Error("%v", err)
		return returnSlice
	}
	defer rows.Close()

	for rows.Next() {
		var job DispatchedJob

		err = rows.Scan(&job.UUID, &job.Runner, &job.Job, &job.Status, &job.LastUpdated)
		if err != nil {
			h.logger.Error("%v", err)
			continue
		}

		result, err := job.ToController()
		if err != nil {
			h.logger.Error("%v", err)
			continue
		}

		returnSlice = append(returnSlice, result)
	}

	return returnSlice
}

type DispatchedJob struct {
	UUID        string
	Runner      string
	Job         []byte
	Status      []byte
	LastUpdated time.Time
}

func (dj *DispatchedJob) ToController() (controller.DispatchedJob, error) {
	result := controller.DispatchedJob{
		UUID:        controller.UUID(dj.UUID),
		Runner:      dj.Runner,
		LastUpdated: dj.LastUpdated,
	}

	if err := json.Unmarshal(dj.Job, &result.Job); err != nil {
		return result, err
	}

	if err := json.Unmarshal(dj.Status, &result.Status); err != nil {
		return result, err
	}

	return result, nil
}

// DeleteJob deletes a specific job from the database.
func (h *HealthCheckerAdapter) DeleteJob(uuid controller.UUID) error {
	_, err := h.db.Client.Exec("DELETE FROM dispatched_jobs WHERE uuid = $1;", uuid)
	return err
}
