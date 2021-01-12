package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/BrenekH/project-redcedar-controller/mediainfo"
)

func main() {
	windowsMediaInfo := "MediaInfo.exe"
	err := mediainfo.SetMediaInfoBinary(windowsMediaInfo)
	if err != nil {
		log.Fatal(err)
	}

	mediainfo, err := mediainfo.GetMediaInfo("I:/input.avi")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mediainfo)
	info, _ := json.Marshal(mediainfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(info))
}
