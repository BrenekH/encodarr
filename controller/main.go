package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BrenekH/project-redcedar-controller/config"
	"github.com/BrenekH/project-redcedar-controller/controller"
	"github.com/BrenekH/project-redcedar-controller/server"
)

func main() {
	wg := &sync.WaitGroup{}
	signals := make(chan os.Signal, 1)
	stopChan := make(chan interface{})
	updateChan := make(chan string)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Printf("Received stop signal: %v", sig)
		stopChan <- true
	}()

	wg.Add(1)
	go controller.RunController(&config.ControllerConfiguration{UpdateChan: &updateChan,
		SearchDir:               "I:/redcedar_test_env",
		FileSystemCheckInterval: 10 * 1e9, // Nanoseconds are stupid
		HealthCheckInterval:     10 * 1e9},
		&stopChan,
		wg)

	wg.Add(1)
	go server.RunHTTPServer(&stopChan, wg)

	<-stopChan

	close(stopChan)
	close(updateChan)

	wg.Wait()
}
