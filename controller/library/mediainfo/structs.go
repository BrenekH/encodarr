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
	// From General Track Type
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

	// From Video Track Type
	ID                             string `json:"ID"`
	IDString                       string `json:"ID_String"`
	FormatInfo                     string `json:"Format_Info"`
	FormatProfile                  string `json:"Format_Profile"`
	FormatLevel                    string `json:"Format_Level"`
	FormatTier                     string `json:"Format_Tier"`
	InternetMediaType              string `json:"InternetMediaType"`
	CodecID                        string `json:"CodecID"`
	Width                          string `json:"Width"`
	WidthString                    string `json:"Width_String"`
	Height                         string `json:"Height"`
	HeightString                   string `json:"Height_String"`
	SampledWidth                   string `json:"Sampled_Width"`
	SampledHeight                  string `json:"Sampled_Height"`
	PixelAspectRatio               string `json:"PixelAspectRatio"`
	DisplayAspectRatio             string `json:"DisplayAspectRatio"`
	DisplayAspectRatioString       string `json:"DisplayAspectRatio_String"`
	FrameRateMode                  string `json:"FrameRate_Mode"`
	FrameRateModeString            string `json:"FrameRate_Mode_String"`
	FrameRateNum                   string `json:"FrameRate_Num"`
	FrameRateDen                   string `json:"FrameRate_Den"`
	ColorSpace                     string `json:"ColorSpace"`
	ChromaSubsampling              string `json:"ChromaSubsampling"`
	ChromaSubsamplingString        string `json:"ChromaSubsampling_String"`
	BitDepth                       string `json:"BitDepth"`
	BitDepthString                 string `json:"BitDepth_String"`
	Delay                          string `json:"Delay"`
	DelayString3                   string `json:"Delay_String3"`
	DelaySource                    string `json:"Delay_Source"`
	DelaySourceString              string `json:"Delay_Source_String"`
	EncodedLibraryName             string `json:"Encoded_Library_Name"`
	EncodedLibraryVersion          string `json:"Encoded_Library_Version"`
	EncodedLibrarySettings         string `json:"Encoded_Library_Settings"`
	Default                        string `json:"Default"`
	DefaultString                  string `json:"Default_String"`
	Forced                         string `json:"Forced"`
	ForcedString                   string `json:"Forced_String"`
	ColourDescriptionPresent       string `json:"colour_description_present"`
	ColourDescriptionPresentSource string `json:"colour_description_present_Source"`
	ColourRange                    string `json:"colour_range"`
	ColourRangeSource              string `json:"colour_range_Source"`
	ColourPrimaries                string `json:"colour_primaries"`
	ColourPrimariesSource          string `json:"colour_primaries_Source"`
	TransferCharacteristics        string `json:"transfer_characteristics"`
	TransferCharacteristicsSource  string `json:"tranfer_characteristics_Source"`
	MatrixCoefficients             string `json:"matrix_coefficients"`
	MatrixCoefficientsSource       string `json:"matrix_coefficients_Source"`

	// From Audio Track Type
	TypeOrder                string `json:"@typeorder"`
	StreamKindPos            string `json:"StreamKindPos"`
	StreamOrder              string `json:"StreamOrder"`
	FormatSettingsSBR        string `json:"Format_Settings_SBR"`
	FormatSettingsSBRString  string `json:"Format_Settings_SBR_String"`
	FormatAdditionalFeatures string `json:"Format_AdditionalFeatures"`
	Channels                 string `json:"Channels"`
	ChannelsString           string `json:"Channels_String"`
	ChannelPositions         string `json:"ChannelPositions"`
	ChannelPositionsString2  string `json:"ChannelPositions_String2"`
	ChannelLayout            string `json:"ChannelLayout"`
	SamplesPerFrame          string `json:"SamplesPerFrame"`
	SamplingRate             string `json:"SamplingRate"`
	SamplingRateString       string `json:"SamplingRate_String"`
	SamplingCount            string `json:"SamplingCount"`
	CompressionMode          string `json:"Compression_Mode"`
	CompressionModeString    string `json:"Compression_Mode_String"`
	DelayString              string `json:"Delay_String"`
	DelayString1             string `json:"Delay_String1"`
	DelayString2             string `json:"Delay_String2"`
	VideoDelay               string `json:"Video_Delay"`
	VideoDelayString         string `json:"Video_Delay_String"`
	VideoDelayString1        string `json:"Video_Delay_String1"`
	VideoDelayString2        string `json:"Video_Delay_String2"`
	VideoDelayString3        string `json:"Video_Delay_String3"`
	Language                 string `json:"Language"`
	LanguageString           string `json:"Language_String"`
	LanguageString1          string `json:"Language_String1"`
	LanguageString2          string `json:"Language_String2"`
	LanguageString3          string `json:"Language_String3"`
	LanguageString4          string `json:"Language_String4"`

	// From Menu Track Type
	ChaptersPosBegin string `json:"Chapters_Pos_Begin"`
	ChaptersPosEnd   string `json:"Chapters_Pos_End"`
	// There are also more items for the extra key, but they have dynamic keys and they don't matter that much anyway.
}
