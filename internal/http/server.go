package http

import (
	"fmt"
	"net/http"
	"time"
)

const (
	readerHeaderTimeout = 5 * time.Second
	defaultTimeout      = 10 * time.Second
)

func NewServer(handler http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           handler,
		ReadHeaderTimeout: readerHeaderTimeout,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
	}
}
