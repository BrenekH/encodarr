package controller

// FileMetadata contains information about a video file.
type FileMetadata struct {
	General        General         `json:"general"`
	VideoTracks    []VideoTrack    `json:"video_tracks"`
	AudioTracks    []AudioTrack    `json:"audio_tracks"`
	SubtitleTracks []SubtitleTrack `json:"subtitle_tracks"`
}

// NOTE: Track type determined by "@type" for MediaInfo and "codec_type" for FFProbe

// General contains the general information about a media file.
type General struct {
	// It looks like any non-string field will have to be parsed
	Duration float32 `json:"duration"`
}

// VideoTrack contains information about a singular video stream in a media file.
type VideoTrack struct {
	Index int64    `json:"index"` // "StreamOrder" (MI), "index" (FF)
	Codec string `json:"codec"` // Either "AVC", "HEVC", etc.
	// Bitrate        int64    `json:"bitrate"`         // "BitRate" (MI), "bit_rate" (FF) // Not implemented for now because I want bitrate per stream, not overall file.
	Width          int64    `json:"width"`           // "Width" (MI), "width" (FF)
	Height         int64    `json:"height"`          // "Height" (MI), "height" (FF)
	ColorPrimaries string `json:"color_primaries"` // "colour_primaries" (MI), "color_primaries" (FF) Will be different based on which MetadataReader is being used (FF gives "bt2020" while MI gives "BT.2020")
}

// AudioTrack contains information about a singular audio stream in a media file.
type AudioTrack struct {
	Index    int64 `json:"index"`    // "StreamOrder" (MI), "index" (FF)
	Channels int64 `json:"channels"` // "Channels" (MI), "channels" (FF)
}

// SubtitleTrack contains information about a singular text stream in a media file.
type SubtitleTrack struct {
	Index    int64    `json:"index"`    // "StreamOrder" (MI), "index" (FF)
	Language string `json:"language"` // "Language" (MI), "tags.language"
}
