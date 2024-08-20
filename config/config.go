package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	InputFilePath   string // Input file path
	OutputFilePath  string // Output file path
	MetricsFilePath string // Metrics file path
	LogLevel        string // Log level
	MaxWorkers      int    // Number of workers to use
	ChunkSize       int    // Number of rows to process in each chunk
	NumberOfRows   int    // Number of lines in the input file
	Version         string // Version of the application
}

func New() *Config {
	_ = godotenv.Load()


	return &Config{
		InputFilePath:   getEnv("INPUT_FILE_PATH", "data/weather_data.csv"),
		OutputFilePath:  getEnv("OUTPUT_FILE_PATH", "data/output.csv"),
		MetricsFilePath: getEnv("METRICS_FILE_PATH", "data/metrics.log"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		MaxWorkers:      getEnvAsInt("MAX_WORKERS", 10),
		ChunkSize:       getEnvAsInt("CHUNK_SIZE", 1000),
		NumberOfRows:   getEnvAsInt("NUMBER_OF_ROWS", 0),
		Version:         getEnv("VERSION", "1.0.0"),
	}
} // New: Create a new Config instance with default values for the configuration options

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error converting %s to int: %v", key, err)
	}
	return intValue
}

// singleton pattern
var instance *Config

func GetInstance() *Config {
	if instance == nil {
		instance = New()
	}
	return instance
}