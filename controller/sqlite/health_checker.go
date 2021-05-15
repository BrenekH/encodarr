package sqlite

import (
	"encoding/json"

	"github.com/BrenekH/encodarr/controller"
)

func NewHealthCheckerAdapater(db *SQLiteDatabase) HealthCheckerAdapter {
	return HealthCheckerAdapter{db: db}
}

// HealthCheckerAdapter satisfies the controller.HealthCheckerDataStorer interface by turning interface
// requests into SQL requests that are passed on to an underlying SQLiteDatabase.
type HealthCheckerAdapter struct {
	db *SQLiteDatabase
}

func (h *HealthCheckerAdapter) DispatchedJobs() []controller.DispatchedJob {
	returnSlice := make([]controller.DispatchedJob, 0)

	rows, err := h.db.Client.Query("SELECT uuid, runner, job, status, last_updated FROM dispatched_jobs;")
	if err != nil {
		// TODO: Log error
		return returnSlice
	}

	for rows.Next() {
		// Variables to scan into
		dj := controller.DispatchedJob{}
		bJ := []byte("") // bytesJob. For intermediate loading into when scanning the rows
		bS := []byte("") // bytesStatus. For intermediate loading into when scanning the rows

		err = rows.Scan(&dj.UUID, &dj.Runner, &bJ, &bS, &dj.LastUpdated)
		if err != nil {
			// TODO: logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bJ, &dj.Job)
		if err != nil {
			// TODO: logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bS, &dj.Status)
		if err != nil {
			// TODO: logger.Error(err.Error())
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
