package sqlite

import (
	"encoding/json"

	"github.com/BrenekH/encodarr/controller"
)

func NewHealthCheckerAdapter(db *SQLiteDatabase, logger controller.Logger) HealthCheckerAdapter {
	return HealthCheckerAdapter{db: db, logger: logger}
}

// HealthCheckerAdapter satisfies the controller.HealthCheckerDataStorer interface by turning interface
// requests into SQL requests that are passed on to an underlying SQLiteDatabase.
type HealthCheckerAdapter struct {
	db     *SQLiteDatabase
	logger controller.Logger
}

func (h *HealthCheckerAdapter) DispatchedJobs() []controller.DispatchedJob {
	returnSlice := make([]controller.DispatchedJob, 0)

	rows, err := h.db.Client.Query("SELECT uuid, runner, job, status, last_updated FROM dispatched_jobs;")
	if err != nil {
		h.logger.Error("%v", err)
		return returnSlice
	}

	for rows.Next() {
		// Variables to scan into
		dj := controller.DispatchedJob{}
		bJ := []byte("") // bytesJob. For intermediate loading into when scanning the rows
		bS := []byte("") // bytesStatus. For intermediate loading into when scanning the rows

		err = rows.Scan(&dj.UUID, &dj.Runner, &bJ, &bS, &dj.LastUpdated)
		if err != nil {
			h.logger.Error("%v", err)
			continue
		}

		err = json.Unmarshal(bJ, &dj.Job)
		if err != nil {
			h.logger.Error("%v", err)
			continue
		}

		err = json.Unmarshal(bS, &dj.Status)
		if err != nil {
			h.logger.Error("%v", err)
			continue
		}

		returnSlice = append(returnSlice, dj)
	}
	rows.Close()

	return returnSlice
}

func (h *HealthCheckerAdapter) DeleteJob(uuid controller.UUID) error {
	_, err := h.db.Client.Exec("DELETE FROM dispatched_jobs WHERE uuid = $1;", uuid)
	return err
}
