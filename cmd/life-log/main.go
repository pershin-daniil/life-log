// Package main initializes and starts the life-log application.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"

	"github.com/pershin-daniil/life-log/internal/server"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	lg := slog.New(slog.NewTextHandler(os.Stderr, nil))

	srv := server.New(lg)

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srv.Run(ctx) })

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("run: %v", err)
	}

	return nil
}
