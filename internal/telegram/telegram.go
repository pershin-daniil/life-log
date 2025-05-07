// Package telegram provides a minimal integration with Telegram Bot API,
// handling basic message receiving and responding functionality.
package telegram

import (
	"log/slog"
)

const module = "telegram"

type Client struct {
	log *slog.Logger
}

func New(log slog.Logger) *Client {
	return &Client{
		log: log.With("module", module),
	}
}
