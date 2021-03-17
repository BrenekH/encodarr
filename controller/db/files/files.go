package files

import (
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/db"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("db/files")
}

type File struct {
	Path    string
	ModTime time.Time
	Queued  bool
}

// All returns a slice of Files that represent the rows in the database
func All() ([]File, error) {
	rows, err := db.Client.Query("SELECT path, modtime, queued FROM files;")
	if err != nil {
		return nil, err
	}

	returnSlice := make([]File, 0)

	for rows.Next() {
		f := File{}

		err = rows.Scan(&f.Path, &f.ModTime, &f.Queued)
		if err != nil {
			logger.Trace(err.Error())
			continue
		}

		returnSlice = append(returnSlice, f)
	}
	if err = rows.Err(); err != nil {
		logger.Trace(err.Error())
		return nil, err
	}

	if err = rows.Close(); err != nil {
		logger.Trace(err.Error())
		return nil, err
	}

	return returnSlice, nil
}

// File "methods"

// Get uses the Path to look up the rest of the information for a File
func (f *File) Get() error {
	err := db.Client.QueryRow("SELECT modtime, queued FROM files WHERE path = $1;", f.Path).Scan(
		&f.ModTime,
		&f.Queued,
	)

	if err != nil {
		logger.Trace(err.Error())
		return err
	}

	return nil
}

// Insert uses the SQL INSERT statement to save the data.
// This means that Insert will fail if the File has already been saved using Insert.
func (f *File) Insert() error {
	_, err := db.Client.Exec("INSERT INTO files (path, modtime, queued) VALUES ($1, $2, $3);",
		f.Path,
		f.ModTime,
		f.Queued,
	)
	if err != nil {
		logger.Trace(err.Error())
		return err
	}

	return nil
}

// Update uses the SQL UPDATE statement to save the data.
// This means that Update will fail if the File hasn't been saved using Insert or it was deleted.
func (f *File) Update() error {
	_, err := db.Client.Exec("UPDATE files SET path=$1, modtime=$2, queued=$3 WHERE path=$1;",
		f.Path,
		f.ModTime,
		f.Queued,
	)
	if err != nil {
		logger.Trace(err.Error())
		return err
	}

	return nil
}

// Upsert uses the SQLite UPSERT paradigm to save the data.
func (f *File) Upsert() error {
	_, err := db.Client.Exec("INSERT INTO files (path, modtime, queued) VALUES ($1, $2, $3) ON CONFLICT(path) DO UPDATE SET path=$1, modtime=$2, queued=$3;",
		f.Path,
		f.ModTime,
		f.Queued,
	)
	if err != nil {
		logger.Trace(err.Error())
		return err
	}

	return nil
}
