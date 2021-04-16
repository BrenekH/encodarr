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

type ApiV1 struct{}

func (a *ApiV1) SendJobComplete(ctx *context.Context, ji runner.JobInfo, cmdR runner.CommandResults) error {
	var request *http.Request
	var err error

	if !cmdR.Failed {
		filename := "output.mkv"

		r, w := io.Pipe()
		writer := multipart.NewWriter(w)

		go func() {
			defer w.Close()
			defer writer.Close()

			file, err := os.Open(filename)
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

		request, err = http.NewRequestWithContext(*ctx, "POST", "http://localhost:8123/api/runner/v1/job/complete", r)
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", writer.FormDataContentType())
	} else {
		request, err = http.NewRequestWithContext(*ctx, "POST", "http://localhost:8123/api/runner/v1/job/complete", &bytes.Buffer{})
		if err != nil {
			return err
		}
	}

	b, err := json.Marshal(historyEntry{
		UUID:   ji.UUID,
		Failed: cmdR.Failed,
		History: history{
			Filename:          ji.File,
			DateTimeCompleted: time.Now(),
			Warnings:          cmdR.Warnings,
			Errors:            cmdR.Errors,
		},
	})
	if err != nil {
		return err
	}

	request.Header.Add("X-Encodarr-History-Entry", string(b))

	response, err := http.DefaultClient.Do(request)
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
	req, err := http.NewRequestWithContext(*ctx, http.MethodGet, "http://localhost:8123/api/runner/v1/job/request", nil)
	if err != nil {
		return runner.JobInfo{}, err
	}

	req.Header.Set("X-Encodarr-Runner-Name", "Develop")

	resp, err := http.DefaultClient.Do(req)
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

	fPath := "input" + path.Ext(jobInfo.Path)

	f, err := os.Create(fPath)
	if err != nil {
		return runner.JobInfo{}, err
	}

	logger.Info(fmt.Sprintf("Received job for %v", jobInfo.Path))

	_, err = io.Copy(f, resp.Body)

	return runner.JobInfo{
		CommandArgs: genFFmpegCmd(fPath, "output.mkv", jobInfo.Parameters),
		UUID:        jobInfo.UUID,
		File:        jobInfo.Path,
		InFile:      fPath,
		OutFile:     "output.mkv",
		MediaInfo:   jobInfo.RawMediaInfo,
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

	req, err := http.NewRequestWithContext(*ctx, http.MethodPost, "http://localhost:8123/api/runner/v1/job/status", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
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
	UUID         string           `json:"uuid"`
	Path         string           `json:"path"`
	Parameters   jobParameters    `json:"parameters"`
	RawMediaInfo runner.MediaInfo `json:"media_info"`
}

// jobParameters represents the actions that need to be taken against a job.
type jobParameters struct {
	Encode bool   `json:"encode"` // true when the file's video stream needs to be encoded
	Stereo bool   `json:"stereo"` // true when the file is missing a stereo audio track
	Codec  string `json:"codec"`  // the ffmpeg compatible video codec
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

// genFFmpegCmd creates the correct ffmpeg arguments for the input/output filenames and the job parameters.
func genFFmpegCmd(inputFname, outputFname string, params jobParameters) []string {
	var s []string
	if params.Stereo && params.Encode {
		s = []string{"-i", inputFname, "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", params.Codec, "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", outputFname}
	} else if params.Stereo {
		s = []string{"-i", inputFname, "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", outputFname}
	} else if params.Encode {
		s = []string{"-i", inputFname, "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", params.Codec, outputFname}
	}
	return s
}
