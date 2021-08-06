package library

import (
	"os"
	"path/filepath"
)

// GetVideoFilesFromDir returns a string slice of video files, found recursively from dirToSearch.
func GetVideoFilesFromDir(dirToSearch string) ([]string, error) {
	allFiles, err := getFilesFromDir(dirToSearch)
	if err != nil {
		return nil, err
	}
	return filterNonVideoExts(allFiles), nil
}

// getFilesFromDir returns all files in a directory.
func getFilesFromDir(dirToSearch string) ([]string, error) {
	cleanSlashedPath := filepath.ToSlash(filepath.Clean(dirToSearch))
	files := make([]string, 0)

	filepath.Walk(cleanSlashedPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, nil
}

// filterNonVideoExts removes any filepath from the provided slice that doesn't end with a known video file extension.
func filterNonVideoExts(toFilter []string) []string {
	// A named return value is not used here because it initializes a nil slice instead of an empty one
	filtered := make([]string, 0)

	for _, i := range toFilter {
		fileExt := filepath.Ext(i)
		if isVideoFileExt(fileExt) {
			filtered = append(filtered, i)
		}
	}
	return filtered
}

// isVideoFileExt returns a bool representing whether or not a file extension is a video file extension.
func isVideoFileExt(a string) bool {
	validExts := []string{".m4v", ".mp4", ".mkv", ".avi", ".mov", ".webm", ".ogg", ".m4p", ".wmv", ".qt"}

	for _, b := range validExts {
		if b == a {
			return true
		}
	}
	return false
}
