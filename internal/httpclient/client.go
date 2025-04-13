// Package httpclient offers a reusable and configurable HTTP client
// for performing outbound HTTP requests with enhanced control over behavior.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const module = "httpclient"

// Client wraps the configuration and behavior for making HTTP requests.
type Client struct {
	log  *slog.Logger
	http *http.Client
}

// New returns a new instance of Client with default settings.
func New(log *slog.Logger) *Client {
	return &Client{
		log:  log.With("module", module),
		http: &http.Client{},
	}
}

// Do sends an HTTP request with the given method, URL, and optional body,
// and decodes the response into respOut.
//
// The request is executed with the provided context.
// respOut should be a pointer to a value into which the response body will be unmarshaled.
//
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) Do(ctx context.Context, method string, url string, body any, respOut any) error {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("do: %v", err)
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return fmt.Errorf("do: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error(err.Error())
		}
	}()

	if resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.log.Error(err.Error())
		}

		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	if respOut != nil {
		return json.NewDecoder(resp.Body).Decode(respOut)
	}

	return nil
}
