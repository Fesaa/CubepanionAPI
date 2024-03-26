package prometheus

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	slog.Info("Starting Prometheus server", "address", "0.0.0.0:8080")

	if e := http.ListenAndServe("0.0.0.0:8080", mux); e != nil {
		slog.Error("Failed to start Prometheus server", "error", e)
	}
}
