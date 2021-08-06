package http

import (
	"os"
	"time"
)

// OsFS is FSer that uses the os package to fulfill the
// interface requirements.
type OsFS struct{}

// Create wraps os.Create
func (o OsFS) Create(name string) (Filer, error) {
	return os.Create(name)
}

// Open wraps os.Open
func (o OsFS) Open(name string) (Filer, error) {
	return os.Open(name)
}

// TimeNow uses time.Now to satisfy the CurrentTimer interface.
type TimeNow struct{}

// Now wraps time.Now
func (t TimeNow) Now() time.Time {
	return time.Now()
}

// FileMetadata contains information about a video file.
type FileMetadata struct {
	General        General         `json:"general"`
	VideoTracks    []VideoTrack    `json:"video_tracks"`
	AudioTracks    []AudioTrack    `json:"audio_tracks"`
	SubtitleTracks []SubtitleTrack `json:"subtitle_tracks"`
}

// General contains the general information about a media file.
type General struct {
	Duration float32 `json:"duration"`
}

// VideoTrack contains information about a singular video stream in a media file.
type VideoTrack struct {
	Index          int    `json:"index"`
	Codec          string `json:"codec"`
	Bitrate        int    `json:"bitrate"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ColorPrimaries string `json:"color_primaries"`
}

// AudioTrack contains information about a singular audio stream in a media file.
type AudioTrack struct {
	Index    int `json:"index"`
	Channels int `json:"channels"`
}

// SubtitleTrack contains information about a singular text stream in a media file.
type SubtitleTrack struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
}
