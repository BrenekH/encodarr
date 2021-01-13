package controller

import (
	"log"
	"os"
	"path/filepath"
)

// GetVideoFilesFromDir returns a string slice of video files, found recursively from dirToSearch.
func GetVideoFilesFromDir(dirToSearch string) []string {
	allFiles, err := getFilesFromDir(dirToSearch)
	if err != nil {
		log.Fatal(err)
	}
	return filterNonVideoExts(allFiles)
}

// getFilesFromDir returns all files in a directory.
func getFilesFromDir(dirToSearch string) (files []string, _ error) {
	cleanSlashedPath := filepath.ToSlash(filepath.Clean(dirToSearch))

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
func filterNonVideoExts(toFilter []string) (filtered []string) {
	for _, i := range toFilter {
		filtered = append(filtered, i)
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
