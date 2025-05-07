// Package main initializes and starts the life-log application.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/pershin-daniil/life-log/internal/config"
	"github.com/pershin-daniil/life-log/internal/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// Создаем контекст, который будет отменен при получении сигнала SIGINT или SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lg := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Load configuration
	cfg, err := config.Load(lg)
	if err != nil {
		return fmt.Errorf("load config: %v", err)
	}

	srv := server.New(lg, cfg)

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srv.Run(ctx) })

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("run: %v", err)
	}

	return nil
}
