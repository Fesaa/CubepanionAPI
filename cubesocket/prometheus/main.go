package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "netty"
	subsystem = ""
)

var (
	constLabels = map[string]string{
		"service": "cubesocket",
	}

	sessionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "sessions_total"),
		Help:        "Total number of sessions opened",
		ConstLabels: constLabels,
	}, []string{})

	sessionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: prometheus.BuildFQName(namespace, subsystem, "session_duration_minutes"),
		Help: "Duration of sessions in minutes",
		Buckets: []float64{
			1, 2, 5, 10, 20, 30, 40, 50, 60, 90, 120, 180, 240, 300, 360, 420, 480, 540, 600,
		},
	}, []string{})

	activeSessions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "active_sessions"),
		Help:        "Number of active sessions",
		ConstLabels: constLabels,
	}, []string{})

	packetsIn = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "packets_in_total"),
		Help:        "Total number of packets received",
		ConstLabels: constLabels,
	}, []string{"packet", "packetId"})

	packetsOut = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "packets_out_total"),
		Help:        "Total number of packets sent",
		ConstLabels: constLabels,
	}, []string{"packet", "packetId"})

	disconnects = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "disconnects_total"),
		Help:        "Total number of disconnects",
		ConstLabels: constLabels,
	}, []string{"reason"})
)
