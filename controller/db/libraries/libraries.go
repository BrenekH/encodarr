package libraries

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/db"
)

// Library represents a singular row in the libraries table
type Library struct {
	ID              int
	Folder          string
	Priority        int
	FsCheckInterval time.Duration
	Pipeline        pluginPipeline
	Queue           Queue
	FileCache       fileCache
	PathMasks       []string
}

type pluginPipeline struct{} // TODO: Implement

type fileCache struct{} // TODO: Complete

// dBLibrary is an interim struct for converting to and from the data types in memory and in the database.
type dBLibrary struct {
	FsCheckInterval string
	Pipeline        []byte
	Queue           []byte
	FileCache       []byte
	PathMasks       []byte
}

var logger logange.Logger

func init() {
	logger = logange.NewLogger("db/libraries")
}

// All returns a slice of Libraries that represent the rows in the database.
func All() ([]Library, error) {
	rows, err := db.Client.Query("SELECT id, folder, priority, fs_check_interval, pipeline, queue, file_cache, path_masks FROM libraries;")
	if err != nil {
		return nil, err
	}
	returnSlice := make([]Library, 0)

	for rows.Next() {
		// Variables to scan into
		l := Library{}
		d := dBLibrary{}

		if err = rows.Scan(&l.ID, &l.Folder, &l.Priority, &d.FsCheckInterval, &d.Pipeline, &d.Queue, &d.FileCache, &d.PathMasks); err != nil {
			logger.Error(err.Error())
			continue
		}

		if err = l.fromDBLibrary(d); err != nil {
			logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, l)
	}
	rows.Close()

	return returnSlice, nil
}

// Library "methods"

// Get uses the UUID to look up the rest of the information for a Library.
func (l *Library) Get() error {
	d := dBLibrary{}

	err := db.Client.QueryRow("SELECT folder, priority, fs_check_interval, pipeline, queue, file_cache, path_masks FROM libraries WHERE id = $1;", l.ID).Scan(
		&l.Folder,
		&l.Priority,
		&d.FsCheckInterval,
		&d.Pipeline,
		&d.Queue,
		&d.FileCache,
		&d.PathMasks,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if err = l.fromDBLibrary(d); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Insert uses the SQL INSERT statement to save the data.
// This means that Insert will fail if the Library has already been saved using Insert.
func (l *Library) Insert() error {
	d, err := l.toDBLibrary()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO libraries (id, folder, priority, fs_check_interval, pipeline, queue, file_cache, path_masks) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		l.ID,
		l.Folder,
		l.Priority,
		d.FsCheckInterval,
		d.Pipeline,
		d.Queue,
		d.FileCache,
		d.PathMasks,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Update uses the SQL UPDATE statement to save the data.
// This means that Update will fail if the Library hasn't been saved using Insert or it was deleted.
func (l *Library) Update() error {
	dbLib, err := l.toDBLibrary()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("UPDATE libraries SET id=$1, folder=$2, priority=$3, fs_check_interval=$4, pipeline=$5, queue=$6, file_cache=$7, path_masks=$8 WHERE id=$1;",
		l.ID,
		l.Folder,
		l.Priority,
		dbLib.FsCheckInterval,
		dbLib.Pipeline,
		dbLib.Queue,
		dbLib.FileCache,
		dbLib.PathMasks,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Delete deletes the corresponding row in the database.
func (l *Library) Delete() error {
	_, err := db.Client.Exec("DELETE FROM libraries WHERE id = $1;", l.ID)
	return err
}

// toDBLibrary returns an instance of dBLibrary with all of the necessary conversions to save data into the database.
func (l Library) toDBLibrary() (d dBLibrary, err error) {
	d.FsCheckInterval = l.FsCheckInterval.String()

	d.Pipeline, err = json.Marshal(l.Pipeline)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d.Queue, err = json.Marshal(l.Queue)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d.FileCache, err = json.Marshal(l.FileCache)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d.PathMasks, err = json.Marshal(l.PathMasks)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

// fromDBLibrary sets the instantiated variables according to the decoded information from the provided dBLibrary.
func (l *Library) fromDBLibrary(d dBLibrary) error {
	var err error
	if d.FsCheckInterval != "" { // This allows FsCheckInterval to not be set in d, while everything still parses correctly.
		l.FsCheckInterval, err = time.ParseDuration(d.FsCheckInterval)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	if err = json.Unmarshal(d.Pipeline, &l.Pipeline); err != nil {
		logger.Error(err.Error())
		return err
	}

	if err = json.Unmarshal(d.Queue, &l.Queue); err != nil {
		logger.Error(err.Error())
		return err
	}

	if err = json.Unmarshal(d.PathMasks, &l.PathMasks); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
