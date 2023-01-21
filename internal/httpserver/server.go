package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/maraero/image-previewer/internal/logger"
)

type Server struct {
	addr   string
	logger logger.Logger
	mux    *http.ServeMux
	srv    *http.Server
}

func New(l logger.Logger) *Server {
	s := &Server{
		addr:   Addr,
		logger: l,
		mux:    http.NewServeMux(),
	}
	s.configureMux()
	return s
}

func (s *Server) Start() error {
	s.srv = &http.Server{
		Addr:         s.addr,
		Handler:      s.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed: %w", err)
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return err
}

func (s *Server) configureMux() {
	s.mux.Handle("/hello", loggerMiddleware(handleHello(s.logger), s.logger))
}
