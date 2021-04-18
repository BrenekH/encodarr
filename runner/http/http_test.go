package http

import (
	"bytes"
	"context"
	"io"
	netHTTP "net/http"
	"testing"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/encodarr/runner/http/mock"
	"github.com/BrenekH/encodarr/runner/options"
)

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

		c := mock.HTTPClient{
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
		apiV1.httpClient = &mock.HTTPClient{
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
