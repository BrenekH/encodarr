package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/runner"
	"github.com/BrenekH/encodarr/runner/cmdrunner"
	"github.com/BrenekH/encodarr/runner/http"
	"github.com/BrenekH/encodarr/runner/options"
	"github.com/BrenekH/logange"
)

var (
	logger logange.Logger
)

func init() {
	logger = logange.NewLogger("main")
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

	cmdRun := cmdrunner.NewCmdRunner()

	apiV1, err := http.NewAPIv1(
		options.TempDir(),
		options.RunnerName(),
		options.ControllerIP(),
		options.ControllerPort(),
	)
	if err != nil {
		logger.Critical(err.Error())
	}

	runner.Run(&ctx, &apiV1, &cmdRun, false)
}
