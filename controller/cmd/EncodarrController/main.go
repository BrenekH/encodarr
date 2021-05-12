package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/controller"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received stop signal: %v", sig)
		// logger.Info(fmt.Sprintf("Received stop signal: %v", sig)) // logange.Logger
		cancel()
	}()

	// TODO: Replace mocks with actual implemented structs
	mockHealthChecker := controller.MockHealthChecker{}
	mockLibraryManager := controller.MockLibraryManager{}
	mockRunnerCommunicator := controller.MockRunnerCommunicator{}
	mockUserInterfacer := controller.MockUserInterfacer{}

	controller.Run(&ctx, &mockHealthChecker, &mockLibraryManager, &mockRunnerCommunicator, &mockUserInterfacer, false)
}
