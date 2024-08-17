package logger

import (
	"testing"
	"time"
)

func TestLogger_TrackPerformance(t *testing.T) {
	log := NewLogger("info")
	log.TrackPerformance(func() {
		time.Sleep(1 * time.Second)
	}, "test")
}
