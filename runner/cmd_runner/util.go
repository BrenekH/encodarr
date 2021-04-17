package cmd_runner

import (
	"fmt"
	"strconv"
	"time"
)

// parseColonTimeToDuration takes a "HH:MM:SS" and converts it to a time.Duration.
// The hour portion does not have to be <= 24.
func parseColonTimeToDuration(s string) (time.Duration, error) {
	var hrs, mins, secs int64

	_, err := fmt.Sscanf(s, "%d:%d:%d", &hrs, &mins, &secs)
	if err != nil {
		return time.Duration(0), err
	}

	// Making everything into time.Durations probably isn't the best option,
	// but there doesn't seem to be a great option to avoid them and still return a time.Duration.
	return time.Duration(
		time.Hour*time.Duration(hrs) +
			time.Minute*time.Duration(mins) +
			time.Second*time.Duration(secs),
	), nil
}

// parseFFmpegLine parses the fps, time, and speed information from a standard FFmpeg statistics line
// and updates the provided pointers if the parsing doesn't return an error.
func parseFFmpegLine(line string, fps *float64, time *string, speed *float64) {
	if pFps, err := extractFps(line); err == nil {
		*fps = pFps
	} else {
		logger.Trace(err.Error())
	}

	if pTime, err := extractTime(line); err == nil {
		*time = pTime
	} else {
		logger.Trace(err.Error())
	}

	if pSpeed, err := extractSpeed(line); err == nil {
		// Not preventing the speed from being zero allows a divide by zero error when calculating stats.
		// Curiously, Go doesn't panic when performing a divide by zero among floats, only integers.
		// This behavior causes the estimated time remaining to appear as a negative number, which makes no sense.
		if pSpeed != 0.0 {
			*speed = pSpeed
		}
	} else {
		logger.Trace(err.Error())
	}
}

func extractFps(line string) (fps float64, err error) {
	fpsReMatch := fpsRe.FindStringSubmatch(line)
	if len(fpsReMatch) > 1 {
		fps, err = strconv.ParseFloat(fpsReMatch[1], 64)
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("extractFps: fps regex returned too few groups (need at least 2)")
		return
	}
	return
}

func extractTime(line string) (time string, err error) {
	timeReMatch := timeRe.FindStringSubmatch(line)
	if len(timeReMatch) > 1 {
		time = timeReMatch[1]
	} else {
		err = fmt.Errorf("extractTime: time regex returned too few groups (need at least 2)")
		return
	}
	return
}

func extractSpeed(line string) (speed float64, err error) {
	speedReMatch := speedRe.FindStringSubmatch(line)
	if len(speedReMatch) > 1 {
		speed, err = strconv.ParseFloat(speedReMatch[1], 64)
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("extractSpeed: speed regex returned too few groups (need at least 2)")
		return
	}
	return
}
