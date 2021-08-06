package sqlite

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

func NewLibraryManagerAdapter(db *SQLiteDatabase, logger controller.Logger) LibraryManagerAdapter {
	return LibraryManagerAdapter{
		db:     db,
		logger: logger,
	}
}

type LibraryManagerAdapter struct {
	db     *SQLiteDatabase
	logger controller.Logger
}

func (l *LibraryManagerAdapter) Libraries() ([]controller.Library, error) {
	rows, err := l.db.Client.Query("SELECT id, folder, priority, fs_check_interval, cmd_decider_settings, queue, path_masks FROM libraries;")
	if err != nil {
		return nil, err
	}
	returnSlice := make([]controller.Library, 0)

	for rows.Next() {
		// Struct to scan into
		d := dbLibrary{}

		if err = rows.Scan(&d.ID, &d.Folder, &d.Priority, &d.FsCheckInterval, &d.CommandDeciderSettings, &d.Queue, &d.PathMasks); err != nil {
			l.logger.Error(err.Error())
			continue
		}

		lib, err := fromDBLibrary(d)
		if err != nil {
			l.logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, lib)
	}
	rows.Close()

	return returnSlice, nil
}

func (l *LibraryManagerAdapter) Library(id int) (controller.Library, error) {
	row := l.db.Client.QueryRow("SELECT id, folder, priority, fs_check_interval, cmd_decider_settings, queue, path_masks FROM libraries WHERE id = $1;", id)

	d := dbLibrary{}

	err := row.Scan(&d.ID, &d.Folder, &d.Priority, &d.FsCheckInterval, &d.CommandDeciderSettings, &d.Queue, &d.PathMasks)
	if err != nil {
		return controller.Library{}, err
	}

	return fromDBLibrary(d)
}

func (l *LibraryManagerAdapter) SaveLibrary(lib controller.Library) error {
	d, err := toDBLibrary(lib)
	if err != nil {
		return err
	}

	_, err = l.db.Client.Exec("INSERT INTO libraries (id, folder, priority, fs_check_interval, cmd_decider_settings, queue, path_masks) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT(id) DO UPDATE SET id=$1, folder=$2, priority=$3, fs_check_interval=$4, cmd_decider_settings=$5, queue=$6, path_masks=$7;",
		d.ID,
		d.Folder,
		d.Priority,
		d.FsCheckInterval,
		d.CommandDeciderSettings,
		d.Queue,
		d.PathMasks,
	)
	if err != nil {
		l.logger.Error(err.Error())
		return err
	}

	return nil
}

// IsPathDispatched loops through the dispatched_jobs table to determine if any jobs with the provided path have already been dispatched.
func (l *LibraryManagerAdapter) IsPathDispatched(path string) (bool, error) {
	rows, err := l.db.Client.Query("SELECT job FROM dispatched_jobs;")
	if err != nil {
		return true, nil
	}
	defer rows.Close()

	for rows.Next() {
		bJob := []byte{}

		if err = rows.Scan(&bJob); err != nil {
			return true, err
		}

		job := controller.Job{}
		err = json.Unmarshal(bJob, &job)
		if err != nil {
			return true, err
		}

		if job.Path == path {
			return true, nil
		}
	}

	return false, nil
}

func (l *LibraryManagerAdapter) PopDispatchedJob(uuid controller.UUID) (controller.DispatchedJob, error) {
	// Get data from table
	row := l.db.Client.QueryRow("SELECT job, status, runner, last_updated FROM dispatched_jobs WHERE uuid = $1", uuid)

	dJob := controller.DispatchedJob{UUID: uuid}
	bJob := []byte{}
	bStatus := []byte{}

	err := row.Scan(
		&bJob,
		&bStatus,
		&dJob.Runner,
		&dJob.LastUpdated,
	)
	if err != nil {
		return dJob, err
	}

	if err = json.Unmarshal(bJob, &dJob.Job); err != nil {
		return dJob, err
	}

	if err = json.Unmarshal(bStatus, &dJob.Status); err != nil {
		return dJob, err
	}

	// Delete data from table
	if _, err = l.db.Client.Exec("DELETE FROM dispatched_jobs WHERE uuid = $1;", uuid); err != nil {
		return dJob, err
	}

	return dJob, nil
}

func (l *LibraryManagerAdapter) PushHistory(h controller.History) error {
	bW, err := json.Marshal(h.Warnings)
	if err != nil {
		return err
	}

	bE, err := json.Marshal(h.Errors)
	if err != nil {
		return err
	}

	_, err = l.db.Client.Exec("INSERT INTO history (time_completed, filename, warnings, errors) VALUES ($1, $2, $3, $4);",
		h.DateTimeCompleted,
		h.Filename,
		bW,
		bE,
	)
	return err
}

// dbLibrary is an interim struct for converting to and from the data types in memory and in the database.
type dbLibrary struct {
	ID                     int
	Folder                 string
	Priority               int
	CommandDeciderSettings string
	FsCheckInterval        string
	Queue                  []byte
	PathMasks              []byte
}

// fromDBLibrary sets the instantiated variables according to the decoded information from the provided dBLibrary.
func fromDBLibrary(d dbLibrary) (controller.Library, error) {
	l := controller.Library{
		ID:                     d.ID,
		Folder:                 d.Folder,
		Priority:               d.Priority,
		CommandDeciderSettings: d.CommandDeciderSettings,
	}

	var err error
	if d.FsCheckInterval != "" { // This allows FsCheckInterval to not be set in d, while everything still parses correctly.
		l.FsCheckInterval, err = time.ParseDuration(d.FsCheckInterval)
		if err != nil {
			return l, err
		}
	}

	if err = json.Unmarshal(d.Queue, &l.Queue); err != nil {
		return l, err
	}

	if err = json.Unmarshal(d.PathMasks, &l.PathMasks); err != nil {
		return l, err
	}

	return l, nil
}

// toDBLibrary returns an instance of dbLibrary with all of the necessary conversions to save data into the database.
func toDBLibrary(lib controller.Library) (d dbLibrary, err error) {
	d.ID = lib.ID
	d.Folder = lib.Folder
	d.Priority = lib.Priority
	d.CommandDeciderSettings = lib.CommandDeciderSettings

	d.FsCheckInterval = lib.FsCheckInterval.String()

	d.Queue, err = json.Marshal(lib.Queue)
	if err != nil {
		return
	}

	d.PathMasks, err = json.Marshal(lib.PathMasks)
	if err != nil {
		return
	}

	return
}
