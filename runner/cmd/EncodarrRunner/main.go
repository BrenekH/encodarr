package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/encodarr/runner/cmd_runner"
	"github.com/BrenekH/encodarr/runner/http"
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
	logger.Info("Starting Encodarr Runner")
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		logger.Info(fmt.Sprintf("Received stop signal: %v", sig))
		cancel()
	}()

	cmdRun := cmd_runner.NewCmdRunner()
	Run(&ctx, &http.ApiV1{}, &cmdRun)
}

// Run runs the basic loop of the Runner
func Run(ctx *context.Context, c runner.Communicator, r runner.CommandRunner) {
	for {
		if IsContextFinished(ctx) {
			break
		}
		// Send new job request
		ji, err := c.SendNewJobRequest(ctx)
		if err != nil {
			logger.Error(err.Error())
		}

		// Start job with request info
		r.Start(ji)

		for !r.Done() {
			// Get status from job
			status := r.Status()

			// Send status to Controller
			err = c.SendStatus(ctx, ji.UUID, status)
			if err != nil {
				logger.Error(err.Error())
			}

			if IsContextFinished(ctx) {
				break
			}
		}

		// Send job complete
		err = c.SendJobComplete(ctx)
		if err != nil {
			logger.Error(err.Error())
		}

		if IsContextFinished(ctx) {
			break
		}
	}
}

type MockCmdRunner struct{}

func (r *MockCmdRunner) Done() bool {
	return true
}

func (r *MockCmdRunner) Start(s string) {}

func (r *MockCmdRunner) Status() {}

func IsContextFinished(ctx *context.Context) bool {
	select {
	case <-(*ctx).Done():
		return true
	default:
		return false
	}
}
