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
	FsCheckInterval time.Duration
	Pipeline        pluginPipeline
	Queue           queue
	FileCache       fileCache
	PathMasks       pathMasks
}

type pluginPipeline struct{} // TODO: Implement

type queue struct{} // TODO: Complete

type fileCache struct{} // TODO: Complete

type pathMasks struct{} // TODO: Complete

var logger logange.Logger

func init() {
	logger = logange.NewLogger("db/libraries")
}

// All returns a slice of Libraries that represent the rows in the database
func All() ([]Library, error) {
	rows, err := db.Client.Query("SELECT id, folder, fs_check_interval, pipeline, queue, file_cache, path_masks FROM libraries;")
	if err != nil {
		return nil, err
	}
	returnSlice := make([]Library, 0)

	for rows.Next() {
		// Variables to scan into
		l := Library{}
		var fsCI string
		bP := []byte("")  // bytesPipeline. For intermediate loading into when scanning the rows
		bQ := []byte("")  // bytesQueue.
		bFC := []byte("") // bytesFileCache.
		bPM := []byte("") // bytesPathMasks.

		err = rows.Scan(&l.ID, &l.Folder, &fsCI, &bP, &bQ, &bFC, &bPM)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		l.FsCheckInterval, err = time.ParseDuration(fsCI)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bP, &l.Pipeline)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bQ, &l.Queue)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bFC, &l.FileCache)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bPM, &l.PathMasks)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, l)
	}
	rows.Close()

	return returnSlice, nil
}

// Library "methods"

// Get uses the UUID to look up the rest of the information for a Library
func (l *Library) Get() error {
	var fsCI string
	bP := []byte("")
	bQ := []byte("")
	bFC := []byte("")
	bPM := []byte("")

	err := db.Client.QueryRow("SELECT folder, fs_check_interval, pipeline, queue, file_cache, path_masks FROM libraries WHERE id = $1;", l.ID).Scan(
		&l.Folder,
		&fsCI,
		&bP,
		&bQ,
		&bFC,
		&bPM,
	)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	l.FsCheckInterval, err = time.ParseDuration(fsCI)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = json.Unmarshal(bP, &l.Pipeline)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = json.Unmarshal(bQ, &l.Queue)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = json.Unmarshal(bPM, &l.PathMasks)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Insert uses the SQL INSERT statement to save the data.
// This means that Insert will fail if the Library has already been saved using Insert.
func (l *Library) Insert() error {
	fsCI := l.FsCheckInterval.String()

	bP, err := json.Marshal(l.Pipeline)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bQ, err := json.Marshal(l.Queue)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bFC, err := json.Marshal(l.FileCache)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bPM, err := json.Marshal(l.PathMasks)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO libraries (id, folder, fs_check_interval, pipeline, queue, file_cache, path_masks) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		l.ID,
		l.Folder,
		fsCI,
		bP,
		bQ,
		bFC,
		bPM,
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
	fsCI := l.FsCheckInterval.String()

	bP, err := json.Marshal(l.Pipeline)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bQ, err := json.Marshal(l.Queue)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bFC, err := json.Marshal(l.FileCache)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bPM, err := json.Marshal(l.PathMasks)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("UPDATE dispatched_jobs SET id=$1, folder=$2, fs_check_interval=$3, pipeline=$4, queue=$5, file_cache=$6, path_masks=$7 WHERE id=$1;",
		l.ID,
		l.Folder,
		fsCI,
		bP,
		bQ,
		bFC,
		bPM,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Delete deletes the corresponding row in the database
func (l *Library) Delete() error {
	_, err := db.Client.Exec("DELETE FROM libraries WHERE id = $1;", l.ID)
	return err
}
