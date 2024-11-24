package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/devv4n/otus-hw/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	httpServer *http.Server
	app        Application
}

type Application interface{}

func NewServer(app Application, cfg *config.Config) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      loggingMiddleware(mux),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	slog.Info("starting server on ", "address", s.httpServer.Addr)

	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return s.Stop(context.Background())
	case err := <-errChan:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("shutting down server...")

	return s.httpServer.Shutdown(ctx)
}
