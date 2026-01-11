package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/takanoriyanagitani/go-http-post-discard"
)

func main() {
	// Command-line flags
	port := flag.Int("port", 10780, "Port to listen on")
	logInterval := flag.Duration("log-interval", 10*time.Second, "Interval for logging statistics")
	flag.Parse()

	// Set up structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize metrics
	metrics := &discard.Metrics{}

	// Start the periodic logger
	metrics.LogPeriodically(logger, *logInterval)

	// Set up the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: metrics.Handler(),
	}

	logger.Info("starting discard server", "address", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
