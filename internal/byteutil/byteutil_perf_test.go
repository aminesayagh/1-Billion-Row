package byteutil

import (
	"testing"
	"onBillion/internal/tracker"
)

func maxInt(a, b int) bool {
	return a < b
}

func maxFloat64(a, b float64) bool {
	return a < b
}

func TestMaxPerformance(t *testing.T) {
	// Setup input data
	aBytes := DefaultByte{1, 2, 3, 4}
	bBytes := DefaultByte{1, 2, 3, 4}
	aInt := int(100)
	bInt := int(200)
	aFloat := float64(100.123)
	bFloat := float64(200.456)

	// Define the functions to test
	f1 := func() {
		Max(aBytes, bBytes)
	}

	f2 := func() {
		maxInt(aInt, bInt)
	}

	f3 := func() {
		maxFloat64(aFloat, bFloat)
	}

	// Number of iterations
	iterations := 1000000

	// Measure performance of MaxBytes vs maxInt
	avgTimeF1, avgTimeF2 := tracker.MeasureExecutionTime(f1, f2, iterations)
	t.Logf("MaxBytes vs maxInt - Avg Time F1: %v, Avg Time F2: %v", avgTimeF1, avgTimeF2)

	// Measure performance of MaxBytes vs maxFloat64
	avgTimeF1, avgTimeF3 := tracker.MeasureExecutionTime(f1, f3, iterations)
	t.Logf("MaxBytes vs maxFloat64 - Avg Time F1: %v, Avg Time F3: %v", avgTimeF1, avgTimeF3)
}
