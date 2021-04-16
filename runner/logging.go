package runner

import (
	"fmt"
	"log"
	"os"

	"github.com/BrenekH/encodarr/runner/options"
	"github.com/BrenekH/logange"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("runner")

	formatter := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	// Setup the root logger to print info
	rootStdoutHandler := logange.NewStdoutHandler()
	rootStdoutHandler.SetFormatter(formatter)
	rootStdoutHandler.SetLevel(options.LogLevel())

	logange.RootLogger.AddHandler(&rootStdoutHandler)

	// Setup a file handler for the root logger
	rootFileHandler, err := logange.NewFileHandler(fmt.Sprintf("%v/runner.log", options.ConfigDir()))
	if err != nil {
		log.Printf("Error creating rootFileHandler: %v", err)
		os.Exit(10)
		return
	}
	rootFileHandler.SetFormatter(formatter)
	rootFileHandler.SetLevel(options.LogLevel())

	logange.RootLogger.AddHandler(&rootFileHandler)
}
