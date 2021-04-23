package http

import (
	"bytes"
	"context"
	"io"
	netHTTP "net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/encodarr/runner/options"
)

func TestSendJobComplete(t *testing.T) {
	apiV1, err := NewApiV1(options.TempDir(), "", "", "")
	if err != nil {
		t.Errorf("Unexpected error creating apiV1: %v", err)
		return
	}
	// Mock out the filesystem and HTTP client so that no external calls are made.
	// If a test needs to evaluate the filesystem or HTTP client, it can assign its own.
	apiV1.fS = &mockFS{}
	apiV1.httpClient = &mockHTTPClient{}

	t.Run("X-Encodarr-History-Entry Set Appropriately", func(t *testing.T) {
		tests := []struct {
			name     string
			expected string
			inJI     runner.JobInfo
			inCR     runner.CommandResults
			inDate   time.Time
		}{
			{
				name:     "Empty",
				expected: `{"uuid":"","failed":false,"history":{"file":"","datetime_completed":"1970-01-01T00:00:00Z","warnings":[],"errors":[]}}`,
				inJI: runner.JobInfo{
					UUID:          "",
					File:          "",
					InFile:        "",
					OutFile:       "",
					CommandArgs:   []string{},
					MediaDuration: 0,
				},
				inCR: runner.CommandResults{
					Failed:         false,
					JobElapsedTime: 0,
					Warnings:       []string{},
					Errors:         []string{},
				},
				inDate: time.Unix(0, 0).UTC(),
			},
			{
				name:     "Populated",
				expected: `{"uuid":"uuid-4","failed":false,"history":{"file":"/tosearch/media/hi.mkv","datetime_completed":"2000-01-01T00:00:00Z","warnings":["Possible corruption"],"errors":[]}}`,
				inJI: runner.JobInfo{
					UUID: "uuid-4",
					File: "/tosearch/media/hi.mkv",
				},
				inCR: runner.CommandResults{
					Failed:   false,
					Warnings: []string{"Possible corruption"},
					Errors:   []string{},
				},
				inDate: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				hC := mockHTTPClient{
					DoResponse: netHTTP.Response{
						StatusCode: 200,
						Body:       io.NopCloser(&bytes.Buffer{}),
					},
				}
				apiV1.httpClient = &hC
				apiV1.currentTime = &mockCurrentTime{time: test.inDate}

				ctx := context.Background()
				outErr := apiV1.SendJobComplete(&ctx, test.inJI, test.inCR)

				if outErr != nil {
					t.Errorf("unexpected error: %v", outErr)
				}

				outHeader := hC.LastRequest.Header.Get("X-Encodarr-History-Entry")
				if outHeader != test.expected {
					t.Errorf("expected %v but got %v", test.expected, outHeader)
				}

				apiV1.httpClient = &mockHTTPClient{}
				apiV1.currentTime = TimeNow{}
			})
		}
	})

	t.Run("Don't Send a File if the Job Failed", func(t *testing.T) {
		hC := mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&bytes.Buffer{}),
			},
		}
		apiV1.httpClient = &hC

		ctx := context.Background()
		outErr := apiV1.SendJobComplete(&ctx, runner.JobInfo{}, runner.CommandResults{Failed: true})

		if outErr != nil {
			t.Errorf("unexpected error: %v", outErr)
		}

		contentType := hC.LastRequest.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/") {
			t.Errorf("expected %v not to start with 'multipart/'", contentType)
		}
	})

	t.Run("HTTPClient.Do is Called", func(t *testing.T) {
		hC := mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&bytes.Buffer{}),
			},
		}
		apiV1.httpClient = &hC

		ctx := context.Background()
		outErr := apiV1.SendJobComplete(&ctx, runner.JobInfo{}, runner.CommandResults{})

		if outErr != nil {
			t.Errorf("unexpected error: %v", outErr)
		}

		if !hC.doCalled {
			t.Errorf("expected HTTPClient.Do to be called, but it wasn't")
		}

		apiV1.httpClient = &mockHTTPClient{}
	})

	t.Run("Handle 409 Status Code", func(t *testing.T) {
		apiV1.httpClient = &mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 409,
				Body:       io.NopCloser(&bytes.Buffer{}),
			},
		}

		ctx := context.Background()
		outErr := apiV1.SendJobComplete(&ctx, runner.JobInfo{}, runner.CommandResults{})

		if outErr != runner.ErrUnresponsive {
			t.Errorf("expected ErrUnresponsive but got %v", outErr)
		}

		apiV1.httpClient = &mockHTTPClient{}
	})
}

func TestSendNewJobRequest(t *testing.T) {
	apiV1, err := NewApiV1("/tmp", "test", "", "")
	if err != nil {
		t.Errorf("Unexpected error creating apiV1: %v", err)
		return
	}
	// To prevent collisions, NewApiV1 creates a random directory inside the provided TempDir.
	// We set the directory manually to circumvent this behavior.
	apiV1.Dir = "/tmp"

	// Mock out the file system and HTTP client so that no system calls are made.
	apiV1.fS = &mockFS{}
	apiV1.httpClient = &mockHTTPClient{}

	t.Run("JobInfo is Properly Derived From the Response Header", func(t *testing.T) {
		tests := []struct {
			name     string
			inStr    string
			expected runner.JobInfo
		}{
			{
				name:  "Filled in (Encode to HEVC)",
				inStr: `{"uuid": "uuid-4", "path": "/media/testFile.mp4", "parameters": {"encode": true, "stereo": false, "codec": "hevc"}, "media_info": {"general": {"duration": "0"}}}`,
				expected: runner.JobInfo{
					UUID:          "uuid-4",
					File:          "/media/testFile.mp4",
					InFile:        "/tmp/input.mp4",
					OutFile:       "/tmp/output.mkv",
					CommandArgs:   []string{"-i", "/tmp/input.mp4", "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", "hevc", "/tmp/output.mkv"},
					MediaDuration: 0,
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				hC := mockHTTPClient{
					DoResponse: netHTTP.Response{
						StatusCode: 200,
						Body:       io.NopCloser(&bytes.Buffer{}),
						Header:     map[string][]string{"X-Encodarr-Job-Info": {test.inStr}},
					},
				}
				apiV1.httpClient = &hC

				ctx := context.Background()
				outJI, outErr := apiV1.SendNewJobRequest(&ctx)

				if outErr != nil {
					t.Errorf("unexpected error: %v", outErr)
				}

				if !reflect.DeepEqual(outJI, test.expected) {
					t.Errorf("expected %v but got %v", test.expected, outJI)
				}
			})
		}
	})

	t.Run("Runner Name is Set", func(t *testing.T) {
		hC := mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&bytes.Buffer{}),
				Header:     map[string][]string{"X-Encodarr-Job-Info": {`{"media_info": {"general": {"duration": "0"}}}`}},
			},
		}
		apiV1.httpClient = &hC

		ctx := context.Background()
		_, outErr := apiV1.SendNewJobRequest(&ctx)

		if outErr != nil {
			t.Errorf("unexpected error: %v", outErr)
		}

		runnerName := hC.LastRequest.Header.Get("X-Encodarr-Runner-Name")
		if runnerName != apiV1.RunnerName {
			t.Errorf("expected %v but got %v", apiV1.RunnerName, runnerName)
		}
	})

	t.Run("Saved File Check", func(t *testing.T) {
		hC := mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&bytes.Buffer{}),
				Header:     map[string][]string{"X-Encodarr-Job-Info": {`{"media_info": {"general": {"duration": "0"}}}`}},
			},
		}
		fS := mockFS{}
		apiV1.httpClient = &hC
		apiV1.fS = &fS

		ctx := context.Background()
		outJI, outErr := apiV1.SendNewJobRequest(&ctx)

		if outErr != nil {
			t.Errorf("unexpected error: %v", outErr)
		}

		if len(fS.createdFiles) < 1 {
			t.Errorf("expected at least one file to created using fS.Create")
		} else {
			lastIndex := len(fS.createdFiles) - 1
			if outJI.InFile != fS.createdFiles[lastIndex] {
				t.Errorf("expected %v but got %v", fS.createdFiles[lastIndex], outJI.InFile)
			}
		}
	})

	t.Run("Ignores 409 Status Code", func(t *testing.T) {
		apiV1.httpClient = &mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 409,
				Body:       io.NopCloser(&bytes.Buffer{}),
				Header:     map[string][]string{"X-Encodarr-Job-Info": {`{"media_info": {"general": {"duration": "0"}}}`}},
			},
		}

		ctx := context.Background()
		_, outErr := apiV1.SendNewJobRequest(&ctx)

		if outErr != nil {
			t.Errorf("expected nil but got %v", outErr)
		}

		apiV1.httpClient = &mockHTTPClient{}
	})
}

func TestApiV1SendStatus(t *testing.T) {
	apiV1, err := NewApiV1(options.TempDir(), "", "", "")
	if err != nil {
		t.Errorf("Unexpected error creating apiV1: %v", err)
		return
	}

	t.Run("Proper HTTP Request is Made", func(t *testing.T) {
		in := runner.JobStatus{
			Stage:                       "running",
			Percentage:                  "76",
			JobElapsedTime:              "1s",
			FPS:                         "9",
			StageElapsedTime:            "1s",
			StageEstimatedTimeRemaining: "1s",
		}
		out := `{"uuid":"uuid-4","status":{"stage":"running","percentage":"76","job_elapsed_time":"1s","fps":"9","stage_elapsed_time":"1s","stage_estimated_time_remaining":"1s"}}`

		c := mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&bytes.Buffer{}),
			},
		}

		apiV1.httpClient = &c

		ctx := context.Background()
		err = apiV1.SendStatus(&ctx, "uuid-4", in)

		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		b, err := io.ReadAll(c.LastRequest.Body)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if string(b) != out {
			t.Errorf("expected %v but got %v", out, string(b))
		}
	})

	t.Run("Respond to Unresponsive Status Code", func(t *testing.T) {
		apiV1.httpClient = &mockHTTPClient{
			DoResponse: netHTTP.Response{
				StatusCode: 409,
				Body:       io.NopCloser(&bytes.Buffer{}),
			},
		}

		ctx := context.Background()
		err = apiV1.SendStatus(&ctx, "uuid-4", runner.JobStatus{})

		if err == nil {
			t.Errorf("Expected Unresponsive error: %v", err)
		} else if err != runner.ErrUnresponsive {
			t.Errorf("Expected Unresponsive error: %v", err)
		}
	})
}
