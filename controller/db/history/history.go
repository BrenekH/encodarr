package history

import (
	"encoding/json"
	"time"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/db"
)

// History represents a row in the history table
type History struct {
	Filename          string    `json:"file"`
	DateTimeCompleted time.Time `json:"datetime_completed"`
	Warnings          []string  `json:"warnings"`
	Errors            []string  `json:"errors"`
}

var logger logange.Logger

func init() {
	logger = logange.NewLogger("db/history")
}

// All returns a slice of Histories that represent the rows in the database
func All() ([]History, error) {
	rows, err := db.Client.Query("SELECT filename, time_completed, warnings, errors FROM history;")
	if err != nil {
		return nil, err
	}

	returnSlice := make([]History, 0)

	for rows.Next() {
		// Variables to scan into
		h := History{Warnings: make([]string, 0), Errors: make([]string, 0)}
		bW := []byte("")
		bE := []byte("")

		err = rows.Scan(&h.Filename, &h.DateTimeCompleted, &bW, &bE)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bW, &h.Warnings)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		err = json.Unmarshal(bE, &h.Errors)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		returnSlice = append(returnSlice, h)
	}

	return returnSlice, nil
}

// History "methods"

// Save inserts the data in History into the database
func (h *History) Save() error {
	bW, err := json.Marshal(h.Warnings)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bE, err := json.Marshal(h.Errors)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = db.Client.Exec("INSERT INTO history (time_completed, filename, warnings, errors) VALUES ($1, $2, $3, $4);",
		h.DateTimeCompleted,
		h.Filename,
		bW,
		bE,
	)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
