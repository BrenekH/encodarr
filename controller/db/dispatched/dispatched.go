package dispatched

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/BrenekH/encodarr/controller/db"
	"github.com/BrenekH/encodarr/controller/mediainfo"
	"github.com/BrenekH/logange"
)

// DJob represents a singular row in the dispatched_jobs table
type DJob struct {
	UUID        string
	Runner      string
	Job         Job
	Status      JobStatus
	LastUpdated time.Time
}

var logger logange.Logger

func init() {
	logger = logange.NewLogger("db/dispatched")
}

// All returns a slice of DJobs that represent the rows in the database
func All() ([]DJob, error) {
	rows, err := db.Client.Query("SELECT uuid, runner, job, status, last_updated FROM dispatched_jobs;")
	if err != nil {
		return nil, err
	}
	returnSlice := make([]DJob, 0)

	for rows.Next() {
		// Variables to scan into
		dj := DJob{}
		bJ := []byte("") // bytesJob. For intermediate loading into when scanning the rows
		bS := []byte("") // bytesStatus. For intermediate loading into when scanning the rows

		err = rows.Scan(&dj.UUID, &dj.Runner, &bJ, &bS, &dj.LastUpdated)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bJ, &dj.Job)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bS, &dj.Status)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, dj)
	}
	rows.Close()

	return returnSlice, nil
}

// PathInDB looks for the provided path in the dispatched jobs
func PathInDB(path string) (bool, error) {
	rows, err := db.Client.Query("SELECT job FROM dispatched_jobs;")
	if err != nil {
		return false, err
	}

	for rows.Next() {
		j := Job{}
		b := []byte("")

		err = rows.Scan(&b)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(b, &j)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		if j.Path == path {
			if err = rows.Close(); err != nil {
				logger.Error(err.Error())
			}
			return true, nil
		}
	}
	rows.Close()

	return false, nil
}

// DJob "methods"

// Get uses the UUID to look up the rest of the information for a DJob
func (d *DJob) Get() error {
	bJ := []byte("")
	bS := []byte("")

	err := db.Client.QueryRow("SELECT runner, job, status, last_updated FROM dispatched_jobs WHERE uuid = $1;", d.UUID).Scan(
		&d.Runner,
		&bJ,
		&bS,
		&d.LastUpdated,
	)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = json.Unmarshal(bJ, &d.Job)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = json.Unmarshal(bS, &d.Status)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Insert uses the SQL INSERT statement to save the data.
// This means that Insert will fail if the dispatched job has already been saved using Insert.
func (d *DJob) Insert() error {
	bJ, err := json.Marshal(d.Job)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bS, err := json.Marshal(d.Status)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO dispatched_jobs (uuid, runner, job, status, last_updated) VALUES ($1, $2, $3, $4, $5);",
		d.UUID,
		d.Runner,
		bJ,
		bS,
		d.LastUpdated,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Update uses the SQL UPDATE statement to save the data.
// This means that Update will fail if the dispatched job hasn't been saved using Insert or it was deleted.
func (d *DJob) Update() error {
	bJ, err := json.Marshal(d.Job)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bS, err := json.Marshal(d.Status)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("UPDATE dispatched_jobs SET uuid=$1, runner=$2, job=$3, status=$4, last_updated=$5 WHERE uuid=$1;",
		d.UUID,
		d.Runner,
		bJ,
		bS,
		d.LastUpdated,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Upsert uses the SQLite UPSERT paradigm to save the data.
func (d *DJob) Upsert() error {
	bJ, err := json.Marshal(d.Job)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bS, err := json.Marshal(d.Status)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO dispatched_jobs (uuid, runner, job, status, last_updated) VALUES ($1, $2, $3, $4, $5) ON CONFLICT(uuid) DO UPDATE SET uuid=$1, runner=$2, job=$3, status=$4, last_updated=$5;",
		d.UUID,
		d.Runner,
		bJ,
		bS,
		d.LastUpdated,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Delete deletes the corresponding row in the database
func (d *DJob) Delete() error {
	_, err := db.Client.Exec("DELETE FROM dispatched_jobs WHERE uuid = $1;", d.UUID)
	return err
}

// Structs stolen from controller.go to avoid cyclic import errors

// Job represents a job in the Encodarr ecosystem.
type Job struct {
	UUID         string              `json:"uuid"`
	Path         string              `json:"path"`
	Parameters   JobParameters       `json:"parameters"`
	RawMediaInfo mediainfo.MediaInfo `json:"media_info"`
}

// JobParameters represents the actions that need to be taken against a job.
type JobParameters struct {
	Encode bool   `json:"encode"` // true when the file's video stream needs to be encoded
	Stereo bool   `json:"stereo"` // true when the file is missing a stereo audio track
	Codec  string `json:"codec"`  // the ffmpeg compatible video codec
}

// Equal is a custom equality check for the Job type
func (j Job) Equal(check Job) bool {
	if j.UUID != check.UUID {
		return false
	}
	if j.Path != check.Path {
		return false
	}
	if !reflect.DeepEqual(j.Parameters, check.Parameters) {
		return false
	}
	return true
}

// EqualPath is a custom equality check for the Job type that only checks the Path parameter
func (j Job) EqualPath(check Job) bool {
	return j.Path == check.Path
}

// EqualUUID is a custom equality check for the Job type that only checks the UUID parameter
func (j Job) EqualUUID(check Job) bool {
	return j.UUID == check.UUID
}

// JobStatus represents the status of a dispatched job
type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}
