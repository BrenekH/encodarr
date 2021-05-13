package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/controller"
	"github.com/BrenekH/encodarr/controller/globals"
	"github.com/BrenekH/encodarr/controller/job_health"
	"github.com/BrenekH/encodarr/controller/sqlite"
)

func main() {
	log.Printf("Starting Encodarr Controller version %v\n", globals.Version)
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received stop signal: %v\n", sig)
		// logger.Info(fmt.Sprintf("Received stop signal: %v", sig)) // logange.Logger
		cancel()
	}()

	sqliteDatabase, err := sqlite.NewSQLiteDatabase(".")
	if err != nil {
		log.Fatal(err)
	}

	hcDBAdapter := sqlite.NewHealthCheckerAdapater(&sqliteDatabase)
	healthChecker := job_health.NewChecker(&hcDBAdapter)

	// TODO: Replace mocks with actual implemented structs
	mockLibraryManager := controller.MockLibraryManager{}
	mockRunnerCommunicator := controller.MockRunnerCommunicator{}
	mockUserInterfacer := controller.MockUserInterfacer{}

	controller.Run(&ctx, &healthChecker, &mockLibraryManager, &mockRunnerCommunicator, &mockUserInterfacer, false)
}
