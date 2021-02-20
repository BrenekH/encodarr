package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	controllerConfig := config.ControllerConfiguration{
		UpdateChan:              &updateChan,
		SearchDir:               "/tosearch",
		JSONDir:                 "/config",
		FileSystemCheckInterval: int(15 * time.Minute),
		HealthCheckInterval:     int(1 * time.Minute),
		HealthCheckTimeout:      int(1 * time.Hour),
	}

	// Start Controller goroutine
	go controller.RunController(&controllerConfig, &stopChan, wg)

	// Start HTTP Server goroutine
	go server.RunHTTPServer(&stopChan, wg)

	<-stopChan

	close(stopChan)
	close(updateChan)

	wg.Wait()
}
