package http

// parseFFmpegCmd takes a string slice and creates a valid parameter list for FFmpeg
func parseFFmpegCmd(inputFname, outputFname string, cmd []string) []string {
	if len(cmd) == 0 {
		return nil
	}

	var s []string = make([]string, len(cmd))

	for i := range cmd {
		if cmd[i] == "ENCODARR_INPUT_FILE" {
			s[i] = inputFname
		} else {
			s[i] = cmd[i]
		}
	}

	return append(s, outputFname)
}
