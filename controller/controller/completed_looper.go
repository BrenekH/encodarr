package controller

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/encodarr/controller/config"
	"github.com/BrenekH/encodarr/controller/db/dispatched"
	"github.com/BrenekH/encodarr/controller/db/history"
)

// HistoryEntry represents an entry for the history collection
type HistoryEntry struct {
	File              string    `json:"file"`
	DateTimeCompleted time.Time `json:"datetime_completed"`
	Warnings          []string  `json:"warnings"`
	Errors            []string  `json:"errors"`
}

// JobCompleteRequest is a struct for representing a job complete request
type JobCompleteRequest struct {
	UUID    string          `json:"uuid"`
	Failed  bool            `json:"failed"`
	History history.History `json:"history"`
	InFile  string          `json:"-"`
}

// completedLooper is a constant loop that spawns goroutines to handle completed files
func completedLooper(completedChan *chan JobCompleteRequest, stopChan *chan interface{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	defer logger.Info("completedLooper stopped")

	for {
		select {
		default:
			select {
			case val := <-*completedChan:
				go completedHandler(val, wg)
			default:
			}
		case <-*stopChan:
			return
		}

		time.Sleep(time.Duration(int64(0.1 * float64(time.Second)))) // Sleep for 0.1 seconds
	}
}

func completedHandler(r JobCompleteRequest, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	// Look up Job information from the database and remove it
	dJob := dispatched.DJob{UUID: r.UUID}
	err := dJob.Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = dJob.Delete()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	filename := dJob.Job.Path

	if config.Global.SmallerFiles {
		ogi, err := os.Stat(filename)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		newI, err := os.Stat(r.InFile)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		if newI.Size() > ogi.Size() {
			logger.Info(fmt.Sprintf("'%v' does not adhere to Smaller Files setting", filename))
			// TODO: Blacklist file in options.ConfigDir()/size_blacklist.json
			return
		}
	}

	// Delete old file from file system
	err = os.Remove(dJob.Job.Path)
	if err != nil {
		failMessage := fmt.Sprintf("Failed to remove file '%v' because of error: %v", dJob.Job.Path, err)
		logger.Error(failMessage)

		// Set filename to a string with an extra encodarr extension
		fnExt := path.Ext(filename)
		i := strings.LastIndex(filename, fnExt)
		fnWoExt := filename[:i] + strings.Replace(filename[i:], fnExt, "", 1)
		filename = fmt.Sprintf("%v.encodarr%v", fnWoExt, fnExt)

		r.History.Warnings = append(r.History.Warnings, failMessage)
	}

	// TODO: Fix saving with the original filename extension instead of the new one

	// Move new file into the old ones place
	err = MoveFile(r.InFile, filename)
	if err != nil {
		failMessage := fmt.Sprintf("Failed to move file '%v' because of error: %v", dJob.Job.Path, err)
		logger.Error(failMessage)

		r.History.Errors = append(r.History.Errors, failMessage)
	}

	err = r.History.Save()
	if err != nil {
		logger.Error(fmt.Sprintf("Error saving history: %v", err.Error()))
	}
}
