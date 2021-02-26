package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/options"
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
	UUID    string       `json:"uuid"`
	Failed  bool         `json:"failed"`
	History HistoryEntry `json:"history"`
	InFile  string       `json:"-"`
}

// HistoryEntries is an instantiated variable of the HistoryContainer type
var HistoryEntries HistoryContainer = HistoryContainer{sync.Mutex{}, make([]HistoryEntry, 0)}

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

	// Look up Job information from DispatchedJobs and remove from DispatchedJobs
	dJob, err := DispatchedJobs.PopByUUID(r.UUID)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not Pop because of invalid UUID '%v': %v", r.UUID, err))
		return
	}

	err = DispatchedJobs.Save()
	if err != nil {
		logger.Error(fmt.Sprintf("Error saving dispatched jobs: %v", err.Error()))
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

		// Set filename to a string with an extra redcedar extension
		fnExt := path.Ext(filename)
		i := strings.LastIndex(filename, fnExt)
		fnWoExt := filename[:i] + strings.Replace(filename[i:], fnExt, "", 1)
		filename = fmt.Sprintf("%v.redcedar%v", fnWoExt, fnExt)

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

	// Add history entry into container
	HistoryEntries.Add(r.History)

	err = HistoryEntries.Save()
	if err != nil {
		logger.Error(fmt.Sprintf("Error saving history: %v", err.Error()))
	}
}

func readHistoryFile() HistoryContainer {
	// Read/unmarshal json from JSONDir/history.json
	f, err := os.Open(fmt.Sprintf("%v/history.json", options.ConfigDir()))
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to open history.json because of error: %v\n", err))
		return HistoryContainer{sync.Mutex{}, make([]HistoryEntry, 0)}
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to read history.json because of error: %v\n", err))
		return HistoryContainer{sync.Mutex{}, make([]HistoryEntry, 0)}
	}

	var readJSON []HistoryEntry
	err = json.Unmarshal(b, &readJSON)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to unmarshal history.json because of error: %v\n", err))
		return HistoryContainer{sync.Mutex{}, make([]HistoryEntry, 0)}
	}

	// Add into HistoryContainer and return
	return HistoryContainer{sync.Mutex{}, readJSON}
}
