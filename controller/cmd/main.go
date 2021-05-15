package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/controller"
	"github.com/BrenekH/encodarr/controller/globals"
	"github.com/BrenekH/encodarr/controller/job_health"
	"github.com/BrenekH/encodarr/controller/library"
	"github.com/BrenekH/encodarr/controller/library/mediainfo"
	"github.com/BrenekH/encodarr/controller/runner_communicator"
	"github.com/BrenekH/encodarr/controller/settings"
	"github.com/BrenekH/encodarr/controller/sqlite"
	"github.com/BrenekH/encodarr/controller/user_interfacer"
	"github.com/BrenekH/logange"
)

func main() {
	configDir := "."

	// Setup main logger
	mainLogger := logange.NewLogger("main")

	formatter := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	// Setup the root logger to print info
	rootStdoutHandler := logange.NewStdoutHandler()
	rootStdoutHandler.SetFormatter(formatter)
	rootStdoutHandler.SetLevel(logange.LevelInfo)

	logange.RootLogger.AddHandler(&rootStdoutHandler)

	// Root logging to a file
	rootFileHandler, err := logange.NewFileHandler(fmt.Sprintf("%v/controller.log", configDir))
	if err != nil {
		log.Printf("Error creating rootFileHandler: %v", err)
		os.Exit(10)
		return
	}
	rootFileHandler.SetFormatter(formatter)
	rootFileHandler.SetLevel(logange.LevelInfo)

	logange.RootLogger.AddHandler(&rootFileHandler)

	mainLogger.Info("Starting Encodarr Controller version %v", globals.Version)
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		mainLogger.Info("Received stop signal: %v", sig)
		cancel()
	}()

	sqliteDatabase, err := sqlite.NewSQLiteDatabase(configDir)
	if err != nil {
		mainLogger.Critical("%v", err)
	}

	settingsStore, err := settings.NewSettingsStore(configDir)
	if err != nil {
		mainLogger.Critical("NewSettingsStore Error: %v", err)
	}

	// --------------- HealthChecker ---------------
	sqliteHCLogger := logange.NewLogger("sqlite.HCA")
	hcDBAdapter := sqlite.NewHealthCheckerAdapter(&sqliteDatabase, &sqliteHCLogger)

	healthCheckerLogger := logange.NewLogger("JobHealth.Checker")
	healthChecker := job_health.NewChecker(&hcDBAdapter, &settingsStore, &healthCheckerLogger)

	// --------------- LibraryManager ---------------
	sqliteLMLogger := logange.NewLogger("sqlite.LMA")
	lmDBAdapter := sqlite.NewLibraryManagerAdapter(&sqliteDatabase, &sqliteLMLogger)

	mediainfoMRLogger := logange.NewLogger("library/mediainfo.MetadataReader")
	metadataReader := mediainfo.NewMetadataReader(&mediainfoMRLogger)

	lmLogger := logange.NewLogger("library.Manager")
	lm := library.NewManager(&lmLogger, &lmDBAdapter, &metadataReader)

	// --------------- RunnerCommunicator ---------------
	rcLogger := logange.NewLogger("runnerCommunicator")
	rc := runner_communicator.NewRunnerHTTPApiV1(&rcLogger)

	// --------------- UserInterfacer ---------------
	uiLogger := logange.NewLogger("userInterfacer")
	ui := user_interfacer.NewWebHTTPApiV1(&uiLogger)

	runLogger := logange.NewLogger("run")
	controller.Run(&ctx, &runLogger, &healthChecker, &lm, &rc, &ui, false)
}
