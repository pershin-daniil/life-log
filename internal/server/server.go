// Package server provides a simple HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

const (
	module          = "server"
	shutdownTimeout = 3 * time.Second
)

// Server wraps an http.Server instance.
type Server struct {
	lg  *slog.Logger
	srv *http.Server
}

// New returns a new Server with default routing.
func New(lg *slog.Logger) *Server {
	version := os.Getenv("VERSION_TAG")
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := fmt.Fprintf(w, "Hello, World!\nversion=%v", version); err != nil {
			return
		}
	})

	return &Server{
		lg: lg.With("module", module),
		srv: &http.Server{
			Addr:              ":8080",
			Handler:           router,
			ReadHeaderTimeout: time.Second,
		},
	}
}

// Run starts the HTTP server and listens for the shutdown signal.
// It returns when the server is gracefully shut down or an error occurs.
func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		gfCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		s.lg.Info("graceful shutdown")
		return s.srv.Shutdown(gfCtx) //nolint:contextcheck // graceful shutdown with new context
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", "addr", s.srv.Addr)

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run: %v", err)
		}

		return nil
	})

	return eg.Wait()
}
