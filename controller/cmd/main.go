package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/controller"
	"github.com/BrenekH/encodarr/controller/cmd/options"
	"github.com/BrenekH/encodarr/controller/globals"
	"github.com/BrenekH/encodarr/controller/httpserver"
	"github.com/BrenekH/encodarr/controller/job_health"
	"github.com/BrenekH/encodarr/controller/library"
	"github.com/BrenekH/encodarr/controller/library/command_decider"
	"github.com/BrenekH/encodarr/controller/library/mediainfo"
	"github.com/BrenekH/encodarr/controller/runner_communicator"
	"github.com/BrenekH/encodarr/controller/settings"
	"github.com/BrenekH/encodarr/controller/sqlite"
	"github.com/BrenekH/encodarr/controller/userinterfacer"
	"github.com/BrenekH/logange"
)

var (
	webApiVersions    = []string{"v1"}
	runnerApiVersions = []string{"v1"}
)

func main() {
	configDir := options.ConfigDir()
	httpServerPort := options.Port()

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

	sqliteDBBuilderLogger := logange.NewLogger("sqlite.DBBuilder")
	sqliteDatabase, err := sqlite.NewSQLiteDatabase(configDir, &sqliteDBBuilderLogger)
	if err != nil {
		mainLogger.Critical("%v", err)
	}

	settingsStore, err := settings.NewSettingsStore(configDir)
	if err != nil {
		mainLogger.Critical("NewSettingsStore Error: %v", err)
	}

	httpSrvLogger := logange.NewLogger("httpServer")
	httpServer := httpserver.NewServer(&httpSrvLogger, httpServerPort, webApiVersions, runnerApiVersions)

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

	fileCacheDS := sqlite.NewFileCacheAdapter(&sqliteDatabase)
	cacheMiddlewareLogger := logange.NewLogger("library.cache")
	metadataCacheMiddleware := library.NewCache(&metadataReader, &fileCacheDS, &cacheMiddlewareLogger)

	cmdDeciderLogger := logange.NewLogger("library/command_decider.CmdDecider")
	commandDecider := command_decider.New(&cmdDeciderLogger)

	lmLogger := logange.NewLogger("library.Manager")
	lm := library.NewManager(&lmLogger, &lmDBAdapter, &metadataCacheMiddleware, &commandDecider)

	// --------------- RunnerCommunicator ---------------
	rcDSLogger := logange.NewLogger("sqlite.RCA")
	rcDS := sqlite.NewRunnerCommunicatorAdapter(&sqliteDatabase, &rcDSLogger)
	rcLogger := logange.NewLogger("runnerCommunicator")
	rc := runner_communicator.NewRunnerHTTPApiV1(&rcLogger, &httpServer, &rcDS)

	// --------------- UserInterfacer ---------------
	uiaLogger := logange.NewLogger("sqlite.UIA")
	uiDBAdapter := sqlite.NewUserInterfacerAdapter(&sqliteDatabase, &uiaLogger)
	uiLogger := logange.NewLogger("userInterfacer")
	ui := userinterfacer.NewWebHTTPv1(&uiLogger, &httpServer, &settingsStore, &uiDBAdapter, false)

	runLogger := logange.NewLogger("run")
	controller.Run(&ctx, &runLogger, &healthChecker, &lm, &rc, &ui, getSetFileLogLevelFunc(&rootFileHandler, &settingsStore), false)
}

func getSetFileLogLevelFunc(fh *logange.FileHandler, ss controller.SettingsStorer) func() {
	return func() {
		switch ss.LogVerbosity() {
		case "TRACE":
			fh.SetLevel(logange.LevelTrace)
		case "DEBUG":
			fh.SetLevel(logange.LevelDebug)
		case "INFO":
			fh.SetLevel(logange.LevelInfo)
		case "WARN":
		case "WARNING":
			fh.SetLevel(logange.LevelWarn)
		case "ERROR":
			fh.SetLevel(logange.LevelError)
		case "CRITICAL":
			fh.SetLevel(logange.LevelCritical)
		default:
			fh.SetLevel(logange.LevelInfo)
		}
	}
}
