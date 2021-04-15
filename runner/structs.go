package runner

import "time"

type JobInfo struct {
	UUID        string
	File        string
	CommandArgs []string
	MediaInfo   MediaInfo
}

type JobStatus struct {
	Stage                       string `json:"stage"`
	Percentage                  string `json:"percentage"`
	JobElapsedTime              string `json:"job_elapsed_time"`
	FPS                         string `json:"fps"`
	StageElapsedTime            string `json:"stage_elapsed_time"`
	StageEstimatedTimeRemaining string `json:"stage_estimated_time_remaining"`
}

type CommandResults struct {
	Failed         bool
	JobElapsedTime time.Duration
	Warnings       []string
	Errors         []string
}

// MediaInfo represents the MediaInfo from a file.
type MediaInfo struct {
	General General `json:"general,omitempty"`
	Video   Video   `json:"video,omitempty"`
	Audio   []Audio `json:"audio,omitempty"`
	Menu    Menu    `json:"menu,omitempty"`
}

type General struct {
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

type Video struct {
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

type Audio struct {
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

type Menu struct {
	Format   string `json:"format"`
	Duration string `json:"duration"`
}
