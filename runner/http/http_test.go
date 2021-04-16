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

func TestApiV1SendStatusUnresponsiveDetection(t *testing.T) {
	apiV1, err := NewApiV1(options.ConfigDir(), "", "", "")
	if err != nil {
		t.Errorf("Unexpected error creating apiV1: %v", err)
	}

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
}
