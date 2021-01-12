package controller

import "github.com/BrenekH/project-redcedar-controller/mediainfo"

// Job represents a job in the RedCedar ecosystem.
type Job struct {
	UUID         string
	Path         string
	Parameters   JobParameters
	RawMediaInfo mediainfo.MediaInfo
}

// JobParameters represents the actions that need to be taken against a job.
type JobParameters struct {
	HEVC        bool // true when the file is not HEVC
	Stereo      bool // true when the file is missing a stereo audio track
	Progressive bool // true when the file is interlaced
}

// RunController is a goroutine compliant way to run the controller.
func RunController() {

}
