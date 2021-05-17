package mediainfo

import "os/exec"

type ExecCommander struct{}

func (e ExecCommander) Command(name string, args ...string) Cmder {
	return exec.Command(name, args...)
}

type mediaInfo struct {
	Media media `json:"media"`
}

type media struct {
	Ref    string  `json:"@ref"`
	Tracks []track `json:"track"`
}

type track struct {
	Type                     string `json:"@type"`
	Count                    string `json:"Count"`
	StreamCount              string `json:"StreamCount"`
	StreamKind               string `json:"StreamKind"`
	StreamKindString         string `json:"StreamKind_String"`
	StreamKindID             string `json:"StreamKindID"`
	UniqueID                 string `json:"UniqueID"`
	UniqueIDString           string `json:"UniqueID_String"`
	VideoCount               string `json:"VideoCount"`
	AudioCount               string `json:"AudioCount"`
	TextCount                string `json:"TextCount"`
	MenuCount                string `json:"MenuCount"`
	VideoFormatList          string `json:"Video_Format_List"`
	VideoFormatWithHintList  string `json:"Video_Format_WithHint_List"`
	VideoCodecList           string `json:"Video_Codec_List"`
	AudioFormatList          string `json:"Audio_Format_List"`
	AudioFormatWithHintList  string `json:"Audio_Format_WithHint_List"`
	AudioCodecList           string `json:"Audio_Codec_List"`
	AudioLanguageList        string `json:"Audio_Language_List"`
	TextFormatList           string `json:"Text_Format_List"`
	TextFormatWithHintList   string `json:"Text_Format_WithHint_List"`
	TextCodecList            string `json:"Text_Codec_List"`
	CompleteName             string `json:"ros.mkv"`
	FileNameExtension        string `json:"FileNameExtension"`
	FileName                 string `json:"FileName"`
	FileExtension            string `json:"FileExtension"`
	Format                   string `json:"Format"`
	FormatString             string `json:"Format_String"`
	FormatUrl                string `json:"Format_Url"`
	FormatExtensions         string `json:"Format_Extensions"`
	FormatCommercial         string `json:"Format_Commercial"`
	FormatVersion            string `json:"Format_Version"`
	FileSize                 string `json:"FileSize"`
	FileSizeString           string `json:"FileSize_String"`
	FileSizeString1          string `json:"FileSize_String1"`
	FileSizeString2          string `json:"FileSize_String2"`
	FileSizeString3          string `json:"FileSize_String3"`
	FileSizeString4          string `json:"FileSize_String4"`
	Duration                 string `json:"Duration"`
	DurationString           string `json:"Duration_String"`
	DurationString1          string `json:"Duration_String1"`
	DurationString2          string `json:"Duration_String2"`
	DurationString3          string `json:"Duration_String3"`
	DurationString4          string `json:"Duration_String4"`
	DurationString5          string `json:"Duration_String5"`
	OverallBitRate           string `json:"OverallBitRate"`
	OverallBitRateString     string `json:"OverallBitRate_String"`
	FrameRate                string `json:"FrameRate"`
	FrameRateString          string `json:"FrameRate_String"`
	FrameCount               string `json:"FrameCount"`
	IsStreamable             string `json:"IsStreamable"`
	Title                    string `json:"Title"`
	Movie                    string `json:"Movie"`
	EncodedDate              string `json:"Encoded_Date"`
	FileCreatedDate          string `json:"File_Created_Date"`
	FileCreatedDateLocal     string `json:"File_Created_Date_Local"`
	FileModifiedDate         string `json:"File_Modified_Date"`
	FileModifiedDateLocal    string `json:"File_Modified_Date_Local"`
	EncodedApplication       string `json:"Encoded_Application"`
	EncodedApplicationString string `json:"Encoded_Appplication_String"`
	EncodedLibrary           string `json:"Encoded_Library"`
	EncodedLibraryString     string `json:"Encoded_Library_String"`
	Extra                    struct {
		ErrorDetectionType string `json:"ErrorDetectionType"`
	} `json:"extra"`
}
