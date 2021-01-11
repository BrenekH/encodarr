package controller

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

// MediainfoBinary specifies a custom MediaInfo binary
var MediainfoBinary = flag.String("mediainfo-bin", "mediainfo", "the path to the mediainfo binary if it is not in the system $PATH")

type mediainfo struct {
	XMLName xml.Name `xml:"MediaInfo"`
	File    file     `xml:"media"`
}

type track struct {
	XMLName                xml.Name `xml:"track"`
	Type                   string   `xml:"type,attr"`
	FileName               string   `xml:"FileName"`
	FormatInfo             string   `xml:"Format_Info"`
	ColorSpace             string   `xml:"ColorSpace"`
	CompleteName           string   `xml:"CompleteName"`
	FormatProfile          string   `xml:"Format_Profile"`
	FileExtension          string   `xml:"FileExtension"`
	ChromaSubsampling      string   `xml:"ChromaSubsampling"`
	WritingApplication     string   `xml:"Writing_application"`
	ProportionOfThisStream string   `xml:"Proportion_of_this_stream"`
	Width                  []string `xml:"Width"`
	Height                 []string `xml:"Height"`
	Format                 []string `xml:"Format"`
	Duration               []string `xml:"Duration"`
	BitRate                []string `xml:"BitRate"`
	BitDepth               []string `xml:"BitDepth"`
	ScanType               []string `xml:"ScanType"`
	FileSize               []string `xml:"FileSize"`
	Framerate              []string `xml:"FrameRate"`
	Channels               []string `xml:"Channels"`
	StreamSize             []string `xml:"StreamSize"`
	BitRateMode            []string `xml:"BitRate_Mode"`
	SamplingRate           []string `xml:"SamplingRate"`
	WritingLibrary         []string `xml:"Writing_library"`
	FramerateMode          []string `xml:"Frame_rate_mode"`
	OverallBitRate         []string `xml:"OverallBitRate"`
	DisplayAspectRatio     []string `xml:"DisplayAspectRatio"`
	OverallBitRateMode     []string `xml:"Overall_bit_rate_mode"`
	FormatSettingsCABAC    []string `xml:"Format_settings__CABAC"`
	FormatSettingsReFrames []string `xml:"Format_settings__ReFrames"`
}

type file struct {
	XMLName xml.Name `xml:"media"`
	Tracks  []track  `xml:"track"`
}

// MediaInfo represents the MediaInfo from a file
type MediaInfo struct {
	General general `json:"general,omitempty"`
	Video   video   `json:"video,omitempty"`
	Audio   audio   `json:"audio,omitempty"`
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
	Framerate          string `json:"frame_rate"`
	StreamSize         string `json:"stream_size"`
	WritingApplication string `json:"writing_application"`
}

type video struct {
	Width                  string `json:"width"`
	Height                 string `json:"height"`
	Format                 string `json:"format"`
	BitRate                string `json:"bitrate"`
	Duration               string `json:"duration"`
	FormatInfo             string `json:"format_info"`
	FormatProfile          string `json:"format_profile"`
	FormatSettingsCABAC    string `json:"format_settings_cabac"`
	FormatSettingsReFrames string `json:"format_settings__reframes"`
	Framerate              string `json:"frame_rate"`
	BitDepth               string `json:"bit_depth"`
	ScanType               string `json:"scan_type"`
	Interlacement          string `json:"interlacement"`
	WritingLibrary         string `json:"writing_library"`
}

type audio struct {
	Format        string `json:"format"`
	Duration      string `json:"duration"`
	BitRate       string `json:"bitrate"`
	Channels      string `json:"channels"`
	Framerate     string `json:"frame_rate"`
	FormatInfo    string `json:"format_Info"`
	SamplingRate  string `json:"sampling_rate"`
	FormatProfile string `json:"format_profile"`
}

type menu struct {
	Format   string `json:"format"`
	Duration string `json:"duration"`
}

// IsInstalled checks if MediaInfo is installed
func IsInstalled() bool {
	cmd := exec.Command(*MediainfoBinary)
	err := cmd.Run()
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") ||
			strings.HasSuffix(err.Error(), "executable file not found in %PATH%") ||
			strings.HasSuffix(err.Error(), "executable file not found in $PATH") {
			return false
		} else if strings.HasPrefix(err.Error(), "exit status 255") {
			return true
		}
	}
	return true
}

// IsMedia checks if the MediaInfo is actual media
func (info MediaInfo) IsMedia() bool {
	return info.Video.Duration != "" && info.Audio.Duration != ""
}

// GetMediaInfo returns MediaInfo from the supplied filename
func GetMediaInfo(fname string) (MediaInfo, error) {
	info := MediaInfo{}
	minfo := mediainfo{}
	general := general{}
	video := video{}
	audio := audio{}
	menu := menu{}

	if !IsInstalled() {
		return info, fmt.Errorf("Must install mediainfo")
	}

	out, err := exec.Command(*MediainfoBinary, "--Output=XML", "-f", fname).Output()

	if err != nil {
		return info, err
	}

	if err := xml.Unmarshal(out, &minfo); err != nil {
		return info, err
	}

	fmt.Println(minfo)

	for _, v := range minfo.File.Tracks {
		if v.Type == "General" {
			general.Duration = v.Duration[0]
			general.Format = v.Format[0]
			general.FileSize = v.FileSize[0]
			if len(v.OverallBitRateMode) > 0 {
				general.OverallBitRateMode = v.OverallBitRateMode[0]
			}
			general.OverallBitRate = v.OverallBitRate[0]
			general.CompleteName = v.CompleteName
			general.FileName = v.FileName
			general.FileExtension = v.FileExtension
			general.Framerate = v.Framerate[0]
			general.StreamSize = v.StreamSize[0]
			general.WritingApplication = v.WritingApplication
		} else if v.Type == "Video" {
			video.Width = v.Width[0]
			video.Height = v.Height[0]
			video.Format = v.Format[0]
			video.BitRate = v.BitRate[0]
			video.Duration = v.Duration[0]
			video.BitDepth = v.BitDepth[0]
			video.ScanType = v.ScanType[0]
			video.FormatInfo = v.FormatInfo
			video.Framerate = v.Framerate[0]
			video.FormatProfile = v.FormatProfile
			if len(v.WritingLibrary) > 0 {
				video.WritingLibrary = v.WritingLibrary[0]
			}
			if len(v.FormatSettingsCABAC) > 0 {
				video.FormatSettingsCABAC = v.FormatSettingsCABAC[0]
			}
			if len(v.FormatSettingsReFrames) > 0 {
				video.FormatSettingsReFrames = v.FormatSettingsReFrames[0]
			}
		} else if v.Type == "Audio" {
			audio.Format = v.Format[0]
			audio.Channels = v.Channels[0]
			audio.Duration = v.Duration[0]
			audio.BitRate = v.BitRate[0]
			audio.FormatInfo = v.FormatInfo
			if len(v.Framerate) > 0 {
				audio.Framerate = v.Framerate[0]
			}
			audio.SamplingRate = v.SamplingRate[0]
			audio.FormatProfile = v.FormatProfile
		} else if v.Type == "Menu" {
			menu.Duration = v.Duration[0]
			menu.Format = v.Format[0]
		}
	}
	info = MediaInfo{General: general, Video: video, Audio: audio, Menu: menu}

	return info, nil
}
