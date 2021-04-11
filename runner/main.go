package main

import (
	"github.com/BrenekH/logange"
)

var (
	logger logange.Logger
)

func init() {
	logger = logange.NewLogger("main")

	// Setup the root logger to print everything
	// This really shouldn't stay here
	f := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	rootStdoutHandler := logange.NewStdoutHandler()
	rootStdoutHandler.SetFormatter(f)
	rootStdoutHandler.SetLevel(logange.LevelTrace)

	logange.RootLogger.AddHandler(&rootStdoutHandler)
}

func main() {
	logger.Info("Hello, World!")
}

// Run runs the basic loop of the Runner
func Run() {
	for {
		// TODO: Send new job request
		// TODO: Start job with request info
		for { // TODO: Add stop condition (job is no longer running)
			// TODO: Get status from job
			// TODO: Send status to Controller
		}
		// TODO: Send job complete
	}
}
