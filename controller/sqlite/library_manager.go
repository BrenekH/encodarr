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

func (l *LibraryManagerAdapter) SaveLibrary(lib controller.Library) error {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return nil
}

func (l *LibraryManagerAdapter) FileEntryByPath(path string) (f controller.File, err error) {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (l *LibraryManagerAdapter) SaveFileEntry(f controller.File) error {
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return nil
}

func (l *LibraryManagerAdapter) IsPathDispatched(path string) (b bool) {
	// Loops through the dispatched_jobs table to determine if any jobs with the provided path have already been dispatched
	l.logger.Critical("Not implemented")
	// TODO: Implement
	return
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
