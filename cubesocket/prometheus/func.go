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

func PacketsIn(packet string, packetId string) {
	packetsIn.WithLabelValues(packet, packetId).Inc()
}

func PacketsOut(packet string, packetId string) {
	packetsOut.WithLabelValues(packet, packetId).Inc()
}

func Disconnect(reason string) {
	disconnects.WithLabelValues(reason).Inc()
}
