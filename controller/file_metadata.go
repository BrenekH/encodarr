package controller

type FileMetadata struct {
	General        General         `json:"general"`
	VideoTracks    []VideoTrack    `json:"video_tracks"`
	AudioTracks    []AudioTrack    `json:"audio_tracks"`
	SubtitleTracks []SubtitleTrack `json:"subtitle_tracks"`
}

// NOTE: Track type determined by "@type" for MediaInfo and "codec_type" for FFProbe

type General struct {
	// It looks like any non-string field will have to be parsed
	Duration float32 `json:"duration"`
}

type VideoTrack struct {
	Index int    `json:"index"` // "StreamOrder" (MI), "index" (FF)
	Codec string `json:"codec"` // Either "AVC", "HEVC", etc.
	// Bitrate        int    `json:"bitrate"`         // "BitRate" (MI), "bit_rate" (FF) // Not implemented for now because I want bitrate per stream, not overall file.
	Width          int    `json:"width"`           // "Width" (MI), "width" (FF)
	Height         int    `json:"height"`          // "Height" (MI), "height" (FF)
	ColorPrimaries string `json:"color_primaries"` // "colour_primaries" (MI), "color_primaries" (FF) Will be different based on which MetadataReader is being used (FF gives "bt2020" while MI gives "BT.2020")
}

type AudioTrack struct {
	Index    int `json:"index"`    // "StreamOrder" (MI), "index" (FF)
	Channels int `json:"channels"` // "Channels" (MI), "channels" (FF)
}

type SubtitleTrack struct {
	Index    int    `json:"index"`    // "StreamOrder" (MI), "index" (FF)
	Language string `json:"language"` // "Language" (MI), "tags.language"
}
