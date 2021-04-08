package config

import (
	"fmt"
	"os"

	"github.com/BrenekH/encodarr/controller/options"
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

	fH, err := logange.NewFileHandler(fmt.Sprintf("%v/controller.log", options.ConfigDir()))
	if err != nil {
		fmt.Printf("Error creating RootFileHandler: %v", err)
		os.Exit(10)
		return
	}
	fH.SetFormatter(f)
	fH.SetLevel(logange.LevelInfo)

	RootFileHandler = fH

	logange.RootLogger.AddHandler(&RootStdoutHandler)
	logange.RootLogger.AddHandler(&RootFileHandler)
}

// SetRootFHVerbosity converts a string into the appropriate logging level and applies it to the root file handler
func SetRootFHVerbosity(v string) error {
	switch v {
	case "TRACE":
		RootFileHandler.SetLevel(logange.LevelTrace)
	case "DEBUG":
		RootFileHandler.SetLevel(logange.LevelDebug)
	case "INFO":
		RootFileHandler.SetLevel(logange.LevelInfo)
	case "WARNING":
		RootFileHandler.SetLevel(logange.LevelWarn)
	case "ERROR":
		RootFileHandler.SetLevel(logange.LevelError)
	case "CRITICAL":
		RootFileHandler.SetLevel(logange.LevelCritical)
	default:
		return fmt.Errorf("invalid logging verbosity '%v'", v)
	}
	return nil
}
