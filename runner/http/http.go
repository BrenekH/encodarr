package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/BrenekH/encodarr/runner"
)

type ApiV1 struct{}

func (a *ApiV1) SendJobComplete(ctx *context.Context) error { return nil }

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

	_, err = io.Copy(f, resp.Body)
	return runner.JobInfo{
		CommandArgs: []string{""},
	}, err
}

func (a *ApiV1) SendStatus(ctx *context.Context) error { return nil }

// Job represents a job in the Encodarr ecosystem.
type Job struct {
	UUID         string        `json:"uuid"`
	Path         string        `json:"path"`
	Parameters   JobParameters `json:"parameters"`
	RawMediaInfo MediaInfo     `json:"media_info"`
}

// JobParameters represents the actions that need to be taken against a job.
type JobParameters struct {
	Encode bool   `json:"encode"` // true when the file's video stream needs to be encoded
	Stereo bool   `json:"stereo"` // true when the file is missing a stereo audio track
	Codec  string `json:"codec"`  // the ffmpeg compatible video codec
}

// MediaInfo represents the MediaInfo from a file.
type MediaInfo struct {
	General general `json:"general,omitempty"`
	Video   video   `json:"video,omitempty"`
	Audio   []audio `json:"audio,omitempty"`
	Menu    menu    `json:"menu,omitempty"`
}

type general struct {
	Format             string `json:"format"`
	Duration           string `json:"duration"`
	FileSize           string `json:"file_size"`
	OverallBitRateMode string `json:"overall_bit_rate_mode"`
	OverallBitRate     string `json:"overall_bit_rate"`
	CompleteName       string `json:"complete_name"`
	FileName           string `json:"file_name"`
	FileExtension      string `json:"file_extension"`
	FrameRate          string `json:"frame_rate"`
	StreamSize         string `json:"stream_size"`
	WritingApplication string `json:"writing_application"`
}

type video struct {
	ID                     string `json:"id"`
	Width                  string `json:"width"`
	Height                 string `json:"height"`
	Format                 string `json:"format"`
	BitRate                string `json:"bitrate"`
	Duration               string `json:"duration"`
	FormatInfo             string `json:"format_info"`
	FormatProfile          string `json:"format_profile"`
	FormatSettingsCABAC    string `json:"format_settings_cabac"`
	FormatSettingsReFrames string `json:"format_settings__reframes"`
	FrameRate              string `json:"frame_rate"`
	BitDepth               string `json:"bit_depth"`
	ScanType               string `json:"scan_type"`
	Interlacement          string `json:"interlacement"`
	WritingLibrary         string `json:"writing_library"`
	ColorPrimaries         string `json:"color_primaries"`
}

type audio struct {
	ID            string `json:"id"`
	Format        string `json:"format"`
	Duration      string `json:"duration"`
	BitRate       string `json:"bitrate"`
	Channels      string `json:"channels"`
	FrameRate     string `json:"frame_rate"`
	FormatInfo    string `json:"format_Info"`
	SamplingRate  string `json:"sampling_rate"`
	FormatProfile string `json:"format_profile"`
}

type menu struct {
	Format   string `json:"format"`
	Duration string `json:"duration"`
}
