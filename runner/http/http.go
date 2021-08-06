package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/logange"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("http")
}

// NewAPIv1 returns an instantiated ApiV1 struct after creating
// a proper temporary directory inside of the provided tempDir argument.
// tempDir will almost always be the result of os.TempDir().
func NewAPIv1(tempDir, runnerName, controllerIP, controllerPort string) (APIv1, error) {
	dir := tempDir + "/Encodarr/Runner"

	if err := os.MkdirAll(dir, 0777); err != nil {
		return APIv1{}, err
	}

	finalDir, err := os.MkdirTemp(dir, "*")
	if err != nil {
		return APIv1{}, err
	}

	return APIv1{
		Dir:          finalDir,
		RunnerName:   runnerName,
		ControllerIP: fmt.Sprintf("http://%v:%v", controllerIP, controllerPort),
		httpClient:   http.DefaultClient,
		fS:           OsFS{},
		currentTime:  TimeNow{},
	}, nil
}

// APIv1 is a struct which implements the runner.Communicator interface using HTTP.
type APIv1 struct {
	Dir          string
	RunnerName   string
	ControllerIP string
	httpClient   RequestDoer
	fS           FSer
	currentTime  CurrentTimer
}

// SendJobComplete lets the Controller know that the job was completed and sends the resulting file if there is one.
func (a *APIv1) SendJobComplete(ctx *context.Context, ji runner.JobInfo, cmdR runner.CommandResults) error {
	var request *http.Request
	var err error

	if !cmdR.Failed {
		filename := a.Dir + "/output.mkv"

		r, w := io.Pipe()
		writer := multipart.NewWriter(w)

		go func() {
			defer w.Close()
			defer writer.Close()

			file, err := a.fS.Open(filename)
			if err != nil {
				logger.Critical(err.Error())
			}
			defer file.Close()

			part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
			if err != nil {
				logger.Critical(err.Error())
			}

			_, err = io.Copy(part, file)
			if err != nil {
				logger.Critical(err.Error())
			}
		}()

		request, err = http.NewRequestWithContext(*ctx, "POST", fmt.Sprintf("%v/api/runner/v1/job/complete", a.ControllerIP), r)
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", writer.FormDataContentType())
	} else {
		request, err = http.NewRequestWithContext(*ctx, "POST", fmt.Sprintf("%v/api/runner/v1/job/complete", a.ControllerIP), &bytes.Buffer{})
		if err != nil {
			return err
		}
	}

	b, err := json.Marshal(historyEntry{
		UUID:   ji.UUID,
		Failed: cmdR.Failed,
		History: history{
			Filename:          ji.File,
			DateTimeCompleted: a.currentTime.Now(),
			Warnings:          cmdR.Warnings,
			Errors:            cmdR.Errors,
		},
	})
	if err != nil {
		return err
	}

	request.Header.Add("X-Encodarr-History-Entry", string(b))

	response, err := a.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 409 {
		return runner.ErrUnresponsive
	}

	return err
}

// SendNewJobRequest requests a new job from the Controller and downloads the file to be worked on.
// This method blocks the thread until a job is assigned to this Runner.
func (a *APIv1) SendNewJobRequest(ctx *context.Context) (runner.JobInfo, error) {
	req, err := http.NewRequestWithContext(*ctx, http.MethodGet, fmt.Sprintf("%v/api/runner/v1/job/request", a.ControllerIP), nil)
	if err != nil {
		return runner.JobInfo{}, err
	}

	req.Header.Set("X-Encodarr-Runner-Name", a.RunnerName)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return runner.JobInfo{}, err
	}
	defer resp.Body.Close()

	strJobInfo := resp.Header.Get("X-Encodarr-Job-Info")

	var jobInfo job
	err = json.Unmarshal([]byte(strJobInfo), &jobInfo)
	if err != nil {
		return runner.JobInfo{}, err
	}

	fPath := a.Dir + "/input" + path.Ext(jobInfo.Path)

	f, err := a.fS.Create(fPath)
	if err != nil {
		return runner.JobInfo{}, err
	}
	defer f.Close()

	logger.Info(fmt.Sprintf("Received job for %v", jobInfo.Path))

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return runner.JobInfo{}, err
	}

	outputFname := a.Dir + "/output.mkv"

	dur := jobInfo.Metadata.General.Duration

	return runner.JobInfo{
		CommandArgs:   parseFFmpegCmd(fPath, outputFname, jobInfo.Command),
		UUID:          jobInfo.UUID,
		File:          jobInfo.Path,
		InFile:        fPath,
		OutFile:       outputFname,
		MediaDuration: dur,
	}, err
}

// SendStatus updates the Controller with the status of the current job.
func (a *APIv1) SendStatus(ctx *context.Context, uuid string, js runner.JobStatus) error {
	b, err := json.Marshal(struct {
		UUID   string           `json:"uuid"`
		Status runner.JobStatus `json:"status"`
	}{
		UUID:   uuid,
		Status: js,
	})
	if err != nil {
		return err
	}

	// TODO: Set a timeout for the http request (using context.Deadline)

	req, err := http.NewRequestWithContext(*ctx, http.MethodPost, fmt.Sprintf("%v/api/runner/v1/job/status", a.ControllerIP), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // We need to close the response body to make sure resources are cleaned up

	if resp.StatusCode == 409 {
		return runner.ErrUnresponsive
	}

	return err
}

// job represents a job in the Encodarr ecosystem.
type job struct {
	UUID     string       `json:"uuid"`
	Path     string       `json:"path"`
	Command  []string     `json:"command"`
	Metadata FileMetadata `json:"metadata"`
}

type historyEntry struct {
	UUID    string  `json:"uuid"`
	Failed  bool    `json:"failed"`
	History history `json:"history"`
}

type history struct {
	Filename          string    `json:"file"`
	DateTimeCompleted time.Time `json:"datetime_completed"`
	Warnings          []string  `json:"warnings"`
	Errors            []string  `json:"errors"`
}
