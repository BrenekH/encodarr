package main

import (
	"context"
	"time"

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
	go func() { time.Sleep(time.Second); cancel() }()
	Run(&ctx, &http.ApiV1{}, &MockCmdRunner{})
}

// Run runs the basic loop of the Runner
func Run(ctx *context.Context, c Communicator, r CommandRunner) {
	for {
		if IsContextFinished(ctx) {
			break
		}
		// TODO: Send new job request
		err := c.SendNewJobRequest(ctx)
		if err != nil {
			logger.Error(err.Error())
		}

		// TODO: Start job with request info
		r.Start()

		for !r.Done() {
			// TODO: Get status from job
			r.Status()

			// TODO: Send status to Controller
			err = c.SendStatus(ctx)
			if err != nil {
				logger.Error(err.Error())
			}

			if IsContextFinished(ctx) {
				break
			}
		}

		// TODO: Send job complete
		err = c.SendJobComplete(ctx)
		if err != nil {
			logger.Error(err.Error())
		}

		if IsContextFinished(ctx) {
			break
		}
	}
}

type Communicator interface {
	SendJobComplete(*context.Context) error
	SendNewJobRequest(*context.Context) error
	SendStatus(*context.Context) error
}

type CommandRunner interface {
	Done() bool
	Start()
	Status()
}

type MockCmdRunner struct{}

func (r *MockCmdRunner) Done() bool {
	return true
}

func (r *MockCmdRunner) Start() {}

func (r *MockCmdRunner) Status() {}

func IsContextFinished(ctx *context.Context) bool {
	select {
	case <-(*ctx).Done():
		return true
	default:
		return false
	}
}
