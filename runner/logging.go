package runner

import "github.com/BrenekH/logange"

var logger logange.Logger

func init() {
	logger = logange.NewLogger("root")

	// Setup the root logger to print everything
	// This really shouldn't stay here
	f := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	rootStdoutHandler := logange.NewStdoutHandler()
	rootStdoutHandler.SetFormatter(f)
	rootStdoutHandler.SetLevel(logange.LevelTrace)

	logange.RootLogger.AddHandler(&rootStdoutHandler)
}
