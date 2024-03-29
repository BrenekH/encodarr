package sqlite

import (
	"encoding/json"

	"github.com/BrenekH/encodarr/controller"
)

// NewUserInterfacerAdapter returns an instantiated UserInterfacerAdapter.
func NewUserInterfacerAdapter(db *Database, logger controller.Logger) UserInterfacerAdapter {
	return UserInterfacerAdapter{db: db, logger: logger}
}

// UserInterfacerAdapter is a struct that satisfies the interface that connects a UserInterfacer
// to a storage medium.
type UserInterfacerAdapter struct {
	db     *Database
	logger controller.Logger
}

// DispatchedJobs returns the content of the dispatched jobs table.
func (u *UserInterfacerAdapter) DispatchedJobs() ([]controller.DispatchedJob, error) {
	returnSlice := make([]controller.DispatchedJob, 0)

	rows, err := u.db.Client.Query("SELECT uuid, runner, job, status, last_updated FROM dispatched_jobs;")
	if err != nil {
		return returnSlice, err
	}

	for rows.Next() {
		// Variables to scan into
		dj := controller.DispatchedJob{}
		bJ := []byte("") // bytesJob. For intermediate loading into when scanning the rows
		bS := []byte("") // bytesStatus. For intermediate loading into when scanning the rows

		err = rows.Scan(&dj.UUID, &dj.Runner, &bJ, &bS, &dj.LastUpdated)
		if err != nil {
			u.logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bJ, &dj.Job)
		if err != nil {
			u.logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bS, &dj.Status)
		if err != nil {
			u.logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, dj)
	}
	rows.Close()

	return returnSlice, nil
}

// HistoryEntries returns the content of the history table.
func (u *UserInterfacerAdapter) HistoryEntries() ([]controller.History, error) {
	returnSlice := make([]controller.History, 0)

	rows, err := u.db.Client.Query("SELECT time_completed, filename, warnings, errors FROM history;")
	if err != nil {
		return returnSlice, err
	}

	for rows.Next() {
		dh := controller.History{}
		bW := []byte("")
		bE := []byte("")

		err = rows.Scan(&dh.DateTimeCompleted, &dh.Filename, &bW, &bE)
		if err != nil {
			u.logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bW, &dh.Warnings)
		if err != nil {
			u.logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bE, &dh.Errors)
		if err != nil {
			u.logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, dh)
	}
	rows.Close()

	return returnSlice, nil
}

// DeleteLibrary deletes the specified library from the libraries table.
func (u *UserInterfacerAdapter) DeleteLibrary(id int) error {
	_, err := u.db.Client.Exec("DELETE FROM libraries WHERE ID = $1;", id)
	return err
}
