package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/maraero/image-previewer/internal/httpserver"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	httpServer := httpserver.New("localhost:3000") // TODO: Move to config
	go func() {
		err := httpServer.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("http server closed: %w", err)
			cancel()
		}
	}()

	fmt.Println("image previewer is running...")
	<-ctx.Done()
	shutdown(httpServer)
}

func shutdown(httpServer *httpserver.Server) {
	fmt.Println("image previewer is turning off...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Stop(ctx); err != nil {
			fmt.Println("failed to stop http server: %w", err.Error())
		} else {
			fmt.Println("http server stopped")
		}
	}()

	wg.Wait()
	fmt.Println("image previewer stopped")
}
