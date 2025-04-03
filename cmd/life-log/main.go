package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context done...")
			return nil
		case t := <-ticker.C:
			slog.Info("msg", "tick", t.Second())
		}
	}
}
