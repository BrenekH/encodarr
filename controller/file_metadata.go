package controller

// TODO: Add json tags (`json:"name"`) to structs

type FileMetadata struct {
	General        General
	VideoTracks    []VideoTrack
	AudioTracks    []AudioTrack
	SubtitleTracks []SubtitleTrack
}

// NOTE: Track type determined by "@type" for MediaInfo and "codec_type" for FFProbe

type General struct {
	// It looks like any non-string field will have to be parsed
	Duration float32
}

type VideoTrack struct {
	Index          int    // "StreamOrder" (MI), "index" (FF)
	Codec          string // Either "AVC", "HEVC", etc.
	Bitrate        int    // "BitRate" (MI), "bit_rate" (FF)
	Width          int    // "Width" (MI), "width" (FF)
	Height         int    // "Height" (MI), "height" (FF)
	ColorPrimaries string // "colour_primaries" (MI), "color_primaries" (FF) Will be different based on which MetadataReader is being used (FF gives "bt2020" while MI gives "BT.2020")
}

type AudioTrack struct {
	Index    int // "StreamOrder" (MI), "index" (FF)
	Channels int // "Channels" (MI), "channels" (FF)
}

type SubtitleTrack struct {
	Index    int    // "StreamOrder" (MI), "index" (FF)
	Language string // "Language" (MI), "tags.language"
}
