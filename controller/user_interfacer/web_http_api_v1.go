package user_interfacer

import (
	"context"
	"sync"

	"github.com/BrenekH/encodarr/controller"
)

func NewWebHTTPApiV1(logger controller.Logger, httpServer controller.HTTPServer) WebHTTPApiV1 {
	return WebHTTPApiV1{logger: logger, httpServer: httpServer}
}

type WebHTTPApiV1 struct {
	logger     controller.Logger
	httpServer controller.HTTPServer
}

func (w *WebHTTPApiV1) Start(ctx *context.Context, wg *sync.WaitGroup) {
	w.httpServer.Start(ctx, wg)

	// TODO: Add handlers to w.httpServer
}

func (w *WebHTTPApiV1) NewLibrarySettings() (m map[int]controller.Library) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
	return
}

func (w *WebHTTPApiV1) SetLibrarySettings([]controller.Library) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
}

func (w *WebHTTPApiV1) SetLibraryQueues([]controller.LibraryQueue) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
}

func (w *WebHTTPApiV1) SetWaitingRunners(runnerNames []string) {
	w.logger.Critical("Not implemented")
	// TODO: Implement
}
