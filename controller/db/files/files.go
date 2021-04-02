package files

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/encodarr/controller/db"
	"github.com/BrenekH/encodarr/controller/mediainfo"
	"github.com/BrenekH/logange"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("db/files")
}

type File struct {
	Path      string
	ModTime   time.Time
	MediaInfo mediainfo.MediaInfo
	Queued    bool
}

// dBFile is an interim struct for converting to and from the data types in memory and in the database.
type dBFile struct {
	MediaInfo []byte
}

// All returns a slice of Files that represent the rows in the database
func All() ([]File, error) {
	rows, err := db.Client.Query("SELECT path, modtime, mediainfo, queued FROM files;")
	if err != nil {
		return nil, err
	}

	returnSlice := make([]File, 0)

	for rows.Next() {
		f := File{}
		d := dBFile{}

		err = rows.Scan(&f.Path, &f.ModTime, &d.MediaInfo, &f.Queued)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		if err = f.fromDBFile(d); err != nil {
			logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, f)
	}
	if err = rows.Err(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if err = rows.Close(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return returnSlice, nil
}

// File "methods"

// Get uses the Path to look up the rest of the information for a File
func (f *File) Get() error {
	d := dBFile{}

	err := db.Client.QueryRow("SELECT modtime, mediainfo, queued FROM files WHERE path = $1;", f.Path).Scan(
		&f.ModTime,
		&d.MediaInfo,
		&f.Queued,
	)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if err = f.fromDBFile(d); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Insert uses the SQL INSERT statement to save the data.
// This means that Insert will fail if the File has already been saved using Insert.
func (f *File) Insert() error {
	d, err := f.toDBFile()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO files (path, modtime, mediainfo, queued) VALUES ($1, $2, $3, $4);",
		f.Path,
		f.ModTime,
		d.MediaInfo,
		f.Queued,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Update uses the SQL UPDATE statement to save the data.
// This means that Update will fail if the File hasn't been saved using Insert or it was deleted.
func (f *File) Update() error {
	d, err := f.toDBFile()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("UPDATE files SET path=$1, modtime=$2, mediainfo=$3, queued=$4 WHERE path=$1;",
		f.Path,
		f.ModTime,
		d.MediaInfo,
		f.Queued,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Upsert uses the SQLite UPSERT paradigm to save the data.
func (f *File) Upsert() error {
	d, err := f.toDBFile()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO files (path, modtime, mediainfo, queued) VALUES ($1, $2, $3, $4) ON CONFLICT(path) DO UPDATE SET path=$1, modtime=$2, mediainfo=$3, queued=$4;",
		f.Path,
		f.ModTime,
		d.MediaInfo,
		f.Queued,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// toDBFile returns an instance of dBFile with all of the necessary conversions to save data into the database.
func (f *File) toDBFile() (dBFile, error) {
	b, err := json.Marshal(f.MediaInfo)
	if err != nil {
		return dBFile{}, err
	}

	return dBFile{MediaInfo: b}, nil
}

// fromDBFile sets the instantiated variables according to the decoded information from the provided dBFile.
func (f *File) fromDBFile(d dBFile) error {
	err := json.Unmarshal(d.MediaInfo, &f.MediaInfo)

	return err
}
