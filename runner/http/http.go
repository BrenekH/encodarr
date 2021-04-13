package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/BrenekH/encodarr/runner"
)

type ApiV1 struct{}

func (a *ApiV1) SendJobComplete(ctx *context.Context) error {
	filename := "output.mkv"

	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

	if err != nil {
		panic(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequestWithContext(*ctx, "POST", "http://localhost:8123/api/runner/v1/job/complete", body)

	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return nil
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

	var jobInfo Job
	err = json.Unmarshal([]byte(strJobInfo), &jobInfo)
	if err != nil {
		return runner.JobInfo{}, err
	}

	fPath := "input" + path.Ext(jobInfo.Path)

	f, err := os.Create(fPath)
	if err != nil {
		return runner.JobInfo{}, err
	}

	fmt.Printf("Received job for %v\n", jobInfo.Path)

	_, err = io.Copy(f, resp.Body)
	return runner.JobInfo{
		CommandArgs: []string{"-i", fPath, "output.mkv"}, // TODO: Construct proper ffmpeg command arguments
		UUID:        jobInfo.UUID,
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

	_, err = http.DefaultClient.Do(req)
	// TODO: Detect if the Controller considers this Runner unresponsive

	return err
}

// Job represents a job in the Encodarr ecosystem.
type Job struct {
	UUID         string           `json:"uuid"`
	Path         string           `json:"path"`
	Parameters   JobParameters    `json:"parameters"`
	RawMediaInfo runner.MediaInfo `json:"media_info"`
}

// JobParameters represents the actions that need to be taken against a job.
type JobParameters struct {
	Encode bool   `json:"encode"` // true when the file's video stream needs to be encoded
	Stereo bool   `json:"stereo"` // true when the file is missing a stereo audio track
	Codec  string `json:"codec"`  // the ffmpeg compatible video codec
}
