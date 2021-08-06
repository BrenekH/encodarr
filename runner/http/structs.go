package http

import (
	"os"
	"time"
)

// OsFS is FSer that uses the os package to fulfill the
// interface requirements.
type OsFS struct{}

func (o OsFS) Create(name string) (Filer, error) {
	return os.Create(name)
}

func (o OsFS) Open(name string) (Filer, error) {
	return os.Open(name)
}

// TimeNow uses time.Now to satisfy the CurrentTimer interface.
type TimeNow struct{}

func (t TimeNow) Now() time.Time {
	return time.Now()
}

type FileMetadata struct {
	General        General         `json:"general"`
	VideoTracks    []VideoTrack    `json:"video_tracks"`
	AudioTracks    []AudioTrack    `json:"audio_tracks"`
	SubtitleTracks []SubtitleTrack `json:"subtitle_tracks"`
}

type General struct {
	Duration float32 `json:"duration"`
}

type VideoTrack struct {
	Index          int    `json:"index"`
	Codec          string `json:"codec"`
	Bitrate        int    `json:"bitrate"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ColorPrimaries string `json:"color_primaries"`
}

type AudioTrack struct {
	Index    int `json:"index"`
	Channels int `json:"channels"`
}

type SubtitleTrack struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
}
