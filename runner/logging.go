package runner

import "github.com/BrenekH/logange"

var logger logange.Logger

func init() {
	logger = logange.NewLogger("runner")

	// Setup the root logger to print info+
	f := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	rootStdoutHandler := logange.NewStdoutHandler()
	rootStdoutHandler.SetFormatter(f)
	rootStdoutHandler.SetLevel(logange.LevelInfo)

	logange.RootLogger.AddHandler(&rootStdoutHandler)

	// TODO: Add file handler for root logger
}
