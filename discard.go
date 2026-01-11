package discard

import (
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

// Metrics holds the performance counters.
type Metrics struct {
	RequestsReceived atomic.Int64
	BytesDiscarded   atomic.Int64
}

// Handler returns an http.Handler that discards the request body and counts metrics.
func (m *Metrics) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.RequestsReceived.Add(1)

		bytes, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			// Optionally log the error, but for a high-performance discard server,
			// we might choose to ignore it to avoid logging overhead per-request.
			// For now, we'll just let the byte count be 0 for this request.
		}
		m.BytesDiscarded.Add(bytes)

		// Respond with 200 OK and no body.
		w.WriteHeader(http.StatusOK)
	})
}

// LogPeriodically starts a goroutine to log metrics at a specified interval.
func (m *Metrics) LogPeriodically(logger *slog.Logger, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			logger.Info("server statistics",
				"reqs_received", m.RequestsReceived.Load(),
				"bytes_discarded", m.BytesDiscarded.Load(),
			)
		}
	}()
}
