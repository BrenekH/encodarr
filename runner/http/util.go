package http

// genFFmpegCmd creates the correct ffmpeg arguments for the input/output filenames and the job parameters.
func genFFmpegCmd(inputFname, outputFname string, params jobParameters) []string {
	var s []string
	if params.Stereo && params.Encode {
		s = []string{"-i", inputFname, "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", params.Codec, "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", outputFname}
	} else if params.Stereo {
		s = []string{"-i", inputFname, "-map", "0:v", "-map", "0:s?", "-map", "0:a", "-map", "0:a", "-c:v", "copy", "-c:s", "copy", "-c:a:1", "copy", "-c:a:0", "aac", "-filter:a:0", "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE", outputFname}
	} else if params.Encode {
		s = []string{"-i", inputFname, "-map", "0:s?", "-map", "0:a", "-c", "copy", "-map", "0:v", "-vcodec", params.Codec, outputFname}
	}

	// Add hardware device if it is not empty
	if params.HWDevice != "" {
		s = append([]string{"-hwaccel_device", params.HWDevice}, s...)
	}

	return s
}
