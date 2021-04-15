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
	runner.Run(&ctx, &http.ApiV1{}, &cmdRun)
}
