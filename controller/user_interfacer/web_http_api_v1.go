package user_interfacer

import (
	"context"
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

func NewWebHTTPApiV1(logger controller.Logger) WebHTTPApiV1 {
	return WebHTTPApiV1{logger: logger}
}

type WebHTTPApiV1 struct {
	logger controller.Logger
}

func (w *WebHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	w.logger.Critical("Not implemented")
	// Run w.httpServer.Start(ctx, wg)
	// Add handlers to w.httpServer
}

func (w *WebHTTPApiV1) NewLibrarySettings() (m map[string]controller.LibrarySettings) {
	w.logger.Critical("Not implemented")
	return
}

func (w *WebHTTPApiV1) SetLibrarySettings([]controller.LibrarySettings) {
	w.logger.Critical("Not implemented")
}

func (w *WebHTTPApiV1) SetLibraryQueues([]controller.LibraryQueue) {
	w.logger.Critical("Not implemented")
}

func (w *WebHTTPApiV1) SetWaitingRunners(runnerNames []string) {
	w.logger.Critical("Not implemented")
}
