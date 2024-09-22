package main

import (
	"testing"
	"oneBillion/config"
	"oneBillion/cmd/version/v2_assembly"
)

// TestMain tests the main function
func TestParsing(t *testing.T) {
	var testConfig = &config.Config{
		InputFilePath:   "data/weather_data_test.csv",
		OutputFilePath:  "data/output_test.csv",
		MetricsFilePath: "data/metrics.log",
		LogLevel:        "info",
		MaxWorkers:      10,
		ChunkSize:       1000,
		NumberOfRows:   0,
		Version:         "0.0.0",
	}

	v2_assembly.Parsing(testConfig)
}