package sse

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEWriter struct {
	writer  http.ResponseWriter
	flusher http.Flusher
}

func SetupSSE(c echo.Context) (*SSEWriter, error) {
	w := c.Response().Writer

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming unsupported")
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	return &SSEWriter{
		writer:  w,
		flusher: flusher,
	}, nil
}

func (s *SSEWriter) SendEvent(event string, data string) error {
	_, err := fmt.Fprintf(s.writer, "event: %s\n", event)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(s.writer, "data: %s\n\n", data)
	if err != nil {
		return err
	}
	s.flusher.Flush()
	return nil
}
