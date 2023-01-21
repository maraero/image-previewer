package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	addr string
	mux  *http.ServeMux
	srv  *http.Server
}

func New() *Server {
	s := &Server{
		addr: Addr,
		mux:  http.NewServeMux(),
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
	s.mux.Handle("/hello", handleHello())
}
