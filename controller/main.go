package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BrenekH/logange"
	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/controller"
	"github.com/BrenekH/project-redcedar-controller/server"
)

var logger logange.Logger

func init() {
	logger = logange.NewLogger("main")
}

func main() {
	defer config.RootFileHandler.Close()

	wg := &sync.WaitGroup{}
	signals := make(chan os.Signal, 1)
	stopChan := make(chan interface{})
	updateChan := make(chan string)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		logger.Info(fmt.Sprintf("Received stop signal: %v", sig))
		stopChan <- true
	}()

	var err error
	config.Global, err = config.LoadSettings()
	if err != nil {
		logger.Warn(err.Error())
		config.Global = config.DefaultSettings()
	}

	err = config.SetRootFHVerbosity(config.Global.LogVerbosity)
	if err != nil {
		logger.Warn(err.Error())
	}

	err = config.SaveGlobal()
	if err != nil {
		logger.Error(err.Error())
	}

	// Start Controller goroutine
	go controller.RunController(&stopChan, wg)

	// Start HTTP Server goroutine
	go server.RunHTTPServer(&stopChan, wg)

	<-stopChan

	close(stopChan)
	close(updateChan)

	wg.Wait()
}
