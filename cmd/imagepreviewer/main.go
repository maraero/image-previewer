package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/maraero/image-previewer/internal/app"
	"github.com/maraero/image-previewer/internal/httpserver"
	"github.com/maraero/image-previewer/internal/imagesrv"
	"github.com/maraero/image-previewer/internal/logger"
)

func main() {
	lggr, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	imageSrv := imagesrv.New(ctx)
	imagepreviewer := app.New(imageSrv)

	httpServer := httpserver.New(imagepreviewer, lggr)
	go func() {
		err := httpServer.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lggr.Fatal("http server closed: %w", err)
			cancel()
		}
	}()

	lggr.Info("image previewer is running...")
	<-ctx.Done()
	shutdown(httpServer, lggr)
}

func shutdown(httpServer *httpserver.Server, logger logger.Logger) {
	logger.Info("image previewer is turning off...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: %w", err.Error())
		} else {
			logger.Info("http server stopped")
		}
	}()

	wg.Wait()
	logger.Info("image previewer stopped")
}
