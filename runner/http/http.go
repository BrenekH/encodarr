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
	"strconv"
	"time"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/logange"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("http")
}

// NewApiV1 returns an instantiated ApiV1 struct after creating
// a proper temporary directory inside of the provided tempDir argument.
// tempDir will almost always be the result of os.TempDir().
func NewApiV1(tempDir, runnerName, controllerIP, controllerPort string) (ApiV1, error) {
	dir := tempDir + "/Encodarr/Runner"

	if err := os.MkdirAll(dir, 0777); err != nil {
		return ApiV1{}, err
	}

	finalDir, err := os.MkdirTemp(dir, "*")
	if err != nil {
		return ApiV1{}, err
	}

	return ApiV1{
		Dir:          finalDir,
		RunnerName:   runnerName,
		ControllerIP: fmt.Sprintf("http://%v:%v", controllerIP, controllerPort),
		httpClient:   http.DefaultClient,
		fS:           OsFS{},
		currentTime:  TimeNow{},
	}, nil
}

// ApiV1 is a struct which implements the runner.Communicator interface using HTTP.
type ApiV1 struct {
	Dir          string
	RunnerName   string
	ControllerIP string
	httpClient   RequestDoer
	fS           FSer
	currentTime  CurrentTimer
}

func (a *ApiV1) SendJobComplete(ctx *context.Context, ji runner.JobInfo, cmdR runner.CommandResults) error {
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

func (a *ApiV1) SendNewJobRequest(ctx *context.Context) (runner.JobInfo, error) {
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

	dur, err := strconv.ParseInt(jobInfo.RawMediaInfo.General.Duration, 10, 64)
	if err != nil {
		return runner.JobInfo{}, err
	}

	return runner.JobInfo{
		CommandArgs:   genFFmpegCmd(fPath, outputFname, jobInfo.Parameters),
		UUID:          jobInfo.UUID,
		File:          jobInfo.Path,
		InFile:        fPath,
		OutFile:       outputFname,
		MediaDuration: dur,
	}, err
}

func (a *ApiV1) SendStatus(ctx *context.Context, uuid string, js runner.JobStatus) error {
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
	UUID         string        `json:"uuid"`
	Path         string        `json:"path"`
	Parameters   jobParameters `json:"parameters"`
	RawMediaInfo MediaInfo     `json:"media_info"`
}

// jobParameters represents the actions that need to be taken against a job.
type jobParameters struct {
	Encode   bool   `json:"encode"`    // true when the file's video stream needs to be encoded
	Stereo   bool   `json:"stereo"`    // true when the file is missing a stereo audio track
	Codec    string `json:"codec"`     // the ffmpeg compatible video codec
	HWDevice string `json:"hw_device"` // The hardware device to use for encoding. If HWDevice is an empty string, a device should not be added to the FFmpeg command.
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
