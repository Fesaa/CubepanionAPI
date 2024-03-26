package prometheus

import "time"

func NewSessions() {
	sessionsTotal.WithLabelValues().Inc()
}

func StartSession() {
	activeSessions.WithLabelValues().Inc()
}

func EndSession() {
	activeSessions.WithLabelValues().Dec()
}

func SessionDuration(duration time.Duration) {
	sessionDuration.WithLabelValues().Observe(duration.Minutes())
}
