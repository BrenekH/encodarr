package mediainfo

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/BrenekH/logange"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("mediainfo")
}

// mediainfoBinary specifies a custom MediaInfo binary.
var mediainfoBinary *string = platformMediaInfoBinary()

func platformMediaInfoBinary() *string {
	var s string
	if runtime.GOOS == "windows" {
		s = "MediaInfo.exe"
	} else {
		s = "mediainfo"
	}
	logger.Debug(fmt.Sprintf("Identified '%v' as the MediaInfo platform binary", s))
	return &s
}

// SetMediaInfoBinary sets the MediaInfo binary to use. Returns an error if the passed binary is invalid.
func SetMediaInfoBinary(s string) error {
	oldVersion := *mediainfoBinary
	mediainfoBinary = &s
	if !IsInstalled() {
		mediainfoBinary = &oldVersion
		return fmt.Errorf("%v is not a valid MediaInfo binary", s)
	}
	return nil
}

type yesNoBoolean bool

func (b *yesNoBoolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s := ""
	d.DecodeElement(&s, &start)

	switch s {
	case "Yes":
		*b = true
	default:
		*b = false
	}

	return nil
}

type mediainfo struct {
	XMLName xml.Name `xml:"Mediainfo"`
	File    file     `xml:"File"`
}

type track struct {
	XMLName                xml.Name     `xml:"track"`
	Type                   string       `xml:"type,attr"`
	Default                yesNoBoolean `xml:"Default"`
	FileName               string       `xml:"File_name"`
	UniqueID               string       `xml:"Unique_ID"`
	FormatInfo             string       `xml:"Format_Info"`
	ColorSpace             string       `xml:"Color_space"`
	CompleteName           string       `xml:"Complete_name"`
	FormatProfile          string       `xml:"Format_profile"`
	FileExtension          string       `xml:"File_extension"`
	ChromaSubsampling      string       `xml:"Chroma_subsampling"`
	WritingApplication     string       `xml:"Writing_application"`
	ProportionOfThisStream string       `xml:"Proportion_of_this_stream"`
	Width                  []string     `xml:"Width"`
	Height                 []string     `xml:"Height"`
	Format                 []string     `xml:"Format"`
	Duration               []string     `xml:"Duration"`
	BitRate                []string     `xml:"Bit_rate"`
	BitDepth               []string     `xml:"Bit_depth"`
	ScanType               []string     `xml:"Scan_type"`
	FileSize               []string     `xml:"File_size"`
	FrameRate              []string     `xml:"Frame_rate"`
	Channels               []string     `xml:"Channel_s_"`
	StreamSize             []string     `xml:"Stream_size"`
	Interlacement          []string     `xml:"Interlacement"`
	BitRateMode            []string     `xml:"Bit_rate_mode"`
	SamplingRate           []string     `xml:"Sampling_rate"`
	WritingLibrary         []string     `xml:"Writing_library"`
	FrameRateMode          []string     `xml:"Frame_rate_mode"`
	OverallBitRate         []string     `xml:"Overall_bit_rate"`
	DisplayAspectRatio     []string     `xml:"Display_aspect_ratio"`
	OverallBitRateMode     []string     `xml:"Overall_bit_rate_mode"`
	FormatSettingsCABAC    []string     `xml:"Format_settings__CABAC"`
	FormatSettingsReFrames []string     `xml:"Format_settings__ReFrames"`
	ColorPrimaries         []string     `xml:"Color_primaries"`
}

type file struct {
	XMLName xml.Name `xml:"File"`
	Tracks  []track  `xml:"track"`
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

// IsInstalled checks if MediaInfo is installed.
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

// IsMedia checks if the MediaInfo is actual media.
func (info MediaInfo) IsMedia() bool {
	return info.Video.Duration != "" && info.Audio[0].Duration != ""
}

func getOrDefault(input []string, index int) string {
	if len(input) > index {
		return input[index]
	}
	return ""
}

// GetMediaInfo returns MediaInfo from the supplied filename.
func GetMediaInfo(fname string) (MediaInfo, error) {
	info := MediaInfo{}
	minfo := mediainfo{}
	mGeneral := general{}
	mVideo := video{}
	mMenu := menu{}
	mAudio := map[string]*audio{}

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

	alreadyParsedVideoStream := false

	for _, v := range minfo.File.Tracks {
		if v.Type == "General" {
			mGeneral.Duration = getOrDefault(v.Duration, 0)
			mGeneral.Format = getOrDefault(v.Format, 0)
			mGeneral.FileSize = getOrDefault(v.FileSize, 0)
			if len(v.OverallBitRateMode) > 0 {
				mGeneral.OverallBitRateMode = getOrDefault(v.OverallBitRateMode, 0)
			}
			mGeneral.OverallBitRate = getOrDefault(v.OverallBitRate, 0)
			mGeneral.CompleteName = v.CompleteName
			mGeneral.FileName = v.FileName
			mGeneral.FileExtension = v.FileExtension
			mGeneral.FrameRate = getOrDefault(v.FrameRate, 0)
			mGeneral.StreamSize = getOrDefault(v.StreamSize, 0)
			mGeneral.WritingApplication = v.WritingApplication
		} else if v.Type == "Video" {
			if alreadyParsedVideoStream && !bool(v.Default) { // Make sure if Default isn't set, we parse at least one track
				// This may not be the most complete way to find parse out video tracks, but we only care about the default video stream for now
				continue
			}
			mVideo.Width = getOrDefault(v.Width, 0)
			mVideo.Height = getOrDefault(v.Height, 0)
			mVideo.Format = getOrDefault(v.Format, 0)
			mVideo.BitRate = getOrDefault(v.BitRate, 0)
			mVideo.Duration = getOrDefault(v.Duration, 0)
			mVideo.BitDepth = getOrDefault(v.BitDepth, 0)
			mVideo.ScanType = getOrDefault(v.ScanType, 0)
			mVideo.FormatInfo = v.FormatInfo
			mVideo.FrameRate = getOrDefault(v.FrameRate, 0)
			mVideo.FormatProfile = v.FormatProfile
			mVideo.Interlacement = getOrDefault(v.Interlacement, 1)
			mVideo.WritingLibrary = getOrDefault(v.WritingLibrary, 0)
			mVideo.FormatSettingsCABAC = getOrDefault(v.FormatSettingsCABAC, 0)
			mVideo.FormatSettingsReFrames = getOrDefault(v.FormatSettingsReFrames, 0)
			mVideo.ColorPrimaries = getOrDefault(v.ColorPrimaries, 0)

			alreadyParsedVideoStream = true
		} else if v.Type == "Audio" {
			audioTrack, inMap := mAudio[v.UniqueID]
			if !inMap {
				mAudio[v.UniqueID] = &audio{}
				audioTrack = mAudio[v.UniqueID]
			}
			audioTrack.Format = getOrDefault(v.Format, 0)
			audioTrack.Channels = getOrDefault(v.Channels, 0)
			audioTrack.Duration = getOrDefault(v.Duration, 0)
			audioTrack.BitRate = getOrDefault(v.BitRate, 0)
			audioTrack.FormatInfo = v.FormatInfo
			audioTrack.FrameRate = getOrDefault(v.FrameRate, 0)
			audioTrack.SamplingRate = getOrDefault(v.SamplingRate, 0)
			audioTrack.FormatProfile = v.FormatProfile
		} else if v.Type == "Menu" {
			mMenu.Duration = getOrDefault(v.Duration, 0)
			mMenu.Format = getOrDefault(v.Format, 0)
		}
	}
	audioSlice := make([]audio, 0)
	for _, v := range mAudio {
		audioSlice = append(audioSlice, *v)
	}
	info = MediaInfo{General: mGeneral, Video: mVideo, Audio: audioSlice, Menu: mMenu}

	return info, nil
}
