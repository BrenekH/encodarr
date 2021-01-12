package mediainfo

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

// mediainfoBinary specifies a custom MediaInfo binary
var mediainfoBinary = flag.String("mediainfo-bin", "mediainfo", "the path to the mediainfo binary if it is not in the system $PATH")

// SetMediaInfoBinary sets the MediaInfo binary to use. Returns an error if the passed binary is invalid
func SetMediaInfoBinary(s string) error {
	oldVersion := *mediainfoBinary
	mediainfoBinary = &s
	if !IsInstalled() {
		mediainfoBinary = &oldVersion
		return fmt.Errorf("%v is not a valid MediaInfo binary", s)
	}
	return nil
}

type mediainfo struct {
	XMLName xml.Name `xml:"Mediainfo"`
	File    file     `xml:"File"`
}

type track struct {
	XMLName                xml.Name `xml:"track"`
	Type                   string   `xml:"type,attr"`
	FileName               string   `xml:"File_name"`
	FormatInfo             string   `xml:"Format_Info"`
	ColorSpace             string   `xml:"Color_space"`
	CompleteName           string   `xml:"Complete_name"`
	FormatProfile          string   `xml:"Format_profile"`
	FileExtension          string   `xml:"File_extension"`
	ChromaSubsampling      string   `xml:"Chroma_subsampling"`
	WritingApplication     string   `xml:"Writing_application"`
	ProportionOfThisStream string   `xml:"Proportion_of_this_stream"`
	Width                  []string `xml:"Width"`
	Height                 []string `xml:"Height"`
	Format                 []string `xml:"Format"`
	Duration               []string `xml:"Duration"`
	BitRate                []string `xml:"Bit_rate"`
	BitDepth               []string `xml:"Bit_depth"`
	ScanType               []string `xml:"Scan_type"`
	FileSize               []string `xml:"File_size"`
	FrameRate              []string `xml:"Frame_rate"`
	Channels               []string `xml:"Channel_s_"`
	StreamSize             []string `xml:"Stream_size"`
	Interlacement          []string `xml:"Interlacement"`
	BitRateMode            []string `xml:"Bit_rate_mode"`
	SamplingRate           []string `xml:"Sampling_rate"`
	WritingLibrary         []string `xml:"Writing_library"`
	FrameRateMode          []string `xml:"Frame_rate_mode"`
	OverallBitRate         []string `xml:"Overall_bit_rate"`
	DisplayAspectRatio     []string `xml:"Display_aspect_ratio"`
	OverallBitRateMode     []string `xml:"Overall_bit_rate_mode"`
	FormatSettingsCABAC    []string `xml:"Format_settings__CABAC"`
	FormatSettingsReFrames []string `xml:"Format_settings__ReFrames"`
}

type file struct {
	XMLName xml.Name `xml:"File"`
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
	FrameRate          string `json:"frame_rate"`
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
	FrameRate              string `json:"frame_rate"`
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
	FrameRate     string `json:"frame_rate"`
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
	cmd := exec.Command(*mediainfoBinary)
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

func getOrDefault(input []string, index int) string {
	if len(input) > index {
		return input[index]
	}
	return ""
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
	out, err := exec.Command(*mediainfoBinary, "--Output=OLDXML", "-f", fname).Output()

	if err != nil {
		return info, err
	}

	if err := xml.Unmarshal(out, &minfo); err != nil {
		return info, err
	}

	for _, v := range minfo.File.Tracks {
		if v.Type == "General" {
			general.Duration = getOrDefault(v.Duration, 0)
			general.Format = getOrDefault(v.Format, 0)
			general.FileSize = getOrDefault(v.FileSize, 0)
			if len(v.OverallBitRateMode) > 0 {
				general.OverallBitRateMode = getOrDefault(v.OverallBitRateMode, 0)
			}
			general.OverallBitRate = getOrDefault(v.OverallBitRate, 0)
			general.CompleteName = v.CompleteName
			general.FileName = v.FileName
			general.FileExtension = v.FileExtension
			general.FrameRate = getOrDefault(v.FrameRate, 0)
			general.StreamSize = getOrDefault(v.StreamSize, 0)
			general.WritingApplication = v.WritingApplication
		} else if v.Type == "Video" {
			video.Width = getOrDefault(v.Width, 0)
			video.Height = getOrDefault(v.Height, 0)
			video.Format = getOrDefault(v.Format, 0)
			video.BitRate = getOrDefault(v.BitRate, 0)
			video.Duration = getOrDefault(v.Duration, 0)
			video.BitDepth = getOrDefault(v.BitDepth, 0)
			video.ScanType = getOrDefault(v.ScanType, 0)
			video.FormatInfo = v.FormatInfo
			video.FrameRate = getOrDefault(v.FrameRate, 0)
			video.FormatProfile = v.FormatProfile
			video.Interlacement = getOrDefault(v.Interlacement, 1)
			video.WritingLibrary = getOrDefault(v.WritingLibrary, 0)
			video.FormatSettingsCABAC = getOrDefault(v.FormatSettingsCABAC, 0)
			video.FormatSettingsReFrames = getOrDefault(v.FormatSettingsReFrames, 0)
		} else if v.Type == "Audio" {
			audio.Format = getOrDefault(v.Format, 0)
			audio.Channels = getOrDefault(v.Channels, 0)
			audio.Duration = getOrDefault(v.Duration, 0)
			audio.BitRate = getOrDefault(v.BitRate, 0)
			audio.FormatInfo = v.FormatInfo
			audio.FrameRate = getOrDefault(v.FrameRate, 0)
			audio.SamplingRate = getOrDefault(v.SamplingRate, 0)
			audio.FormatProfile = v.FormatProfile
		} else if v.Type == "Menu" {
			menu.Duration = getOrDefault(v.Duration, 0)
			menu.Format = getOrDefault(v.Format, 0)
		}
	}
	info = MediaInfo{General: general, Video: video, Audio: audio, Menu: menu}

	return info, nil
}
