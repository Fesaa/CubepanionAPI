package prometheus

import (
	"fmt"
	"time"
)

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

func VersionInc(version int) {
	versions.WithLabelValues(fmt.Sprintf("%d", version)).Inc()
}

func VersionDec(version int) {
	versions.WithLabelValues(fmt.Sprintf("%d", version)).Dec()
}
