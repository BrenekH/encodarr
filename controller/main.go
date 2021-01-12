package main

import (
	"os"
)

// IsDirectory returns a bool representing whether or not the provided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func main() {
	// windowsMediaInfo := "MediaInfo.exe"
	// err := mediainfo.SetMediaInfoBinary(windowsMediaInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// mediainfo, err := mediainfo.GetMediaInfo("I:/test_input.avi")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(mediainfo)
	// info, _ := json.Marshal(mediainfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(info))

	// paths, err := filepath.Glob("I:/redcedar_test_env/**/*")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// p, err := filepath.Glob("I:/redcedar_test_env/*")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// paths = append(paths, p...)

	// filteredPaths := make([]string, 0)

	// for _, path := range paths {
	// 	if b, _ := IsDirectory(path); !b {
	// 		fmt.Println(path)
	// 		filteredPaths = append(filteredPaths, path)
	// 	}
	// }
	// fmt.Println(filteredPaths)

	// fmt.Println(filepath.ToSlash(filepath.Clean("I:/redcedar_test_env/")) + "/**/*")
	// fmt.Println(filepath.ToSlash(filepath.Clean("I:/redcedar_test_env")) + "/*")
}
