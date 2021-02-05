package config

import (
	"fmt"
	"os"

	"github.com/BrenekH/logange"
)

var (
	// RootStdoutHandler is a "globalized" StdoutHandler for the RootLogger
	RootStdoutHandler logange.StdoutHandler
	// RootFileHandler is a "globalized" FileHandler for the RootLogger
	RootFileHandler logange.FileHandler
)

func init() {
	f := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	RootStdoutHandler = logange.NewStdoutHandler()
	RootStdoutHandler.SetFormatter(f)
	RootStdoutHandler.SetLevel(logange.LevelWarn)

	fH, err := logange.NewFileHandler("/config/controller.log")
	if err != nil {
		fmt.Printf("Error creating RootFileHandler: %v", err)
		os.Exit(10)
		return
	}
	fH.SetFormatter(f)
	fH.SetLevel(logange.LevelTrace)

	RootFileHandler = fH

	logange.RootLogger.AddHandler(&RootStdoutHandler)
	logange.RootLogger.AddHandler(&RootFileHandler)
}
