package logger

import (
	"onBillion/config" // what I
	"os"
	"strconv" // convert the float64 values to strings
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerInterface interface {
	TrackPerformance(f func())
}

type Logger struct {
	*logrus.Logger
	mu           sync.Mutex // mutual exclusion lock
	fastestChunk float64
	slowestChunk float64
	metric_name  string
}

func NewLogger(logLevel string) *Logger {
	log := logrus.New()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	return &Logger{
		Logger:       log,
		fastestChunk: 0,
		slowestChunk: 0,
	}
}

// TrackPerformance measures the time taken by a function to execute and logs it
func (l *Logger) TrackPerformance(f func(), metricName string) {
	start := time.Now()
	f()
	elapsed := time.Since(start).Seconds()

	l.mu.Lock()         // lock the mutex before updating the metrics
	defer l.mu.Unlock() // unlock the mutex after the function returns

	if elapsed < l.fastestChunk {
		l.fastestChunk = elapsed
	}

	if elapsed > l.slowestChunk {
		l.slowestChunk = elapsed
	}

	l.metric_name = metricName

	l.WithFields(logrus.Fields{
		"fastest_chunk": l.fastestChunk,
		"slowest_chunk": l.slowestChunk,
		"metric_name":   metricName,
	}).Info("Performance metrics")

	// save the metrics to a file
	l.saveMetrics()
}

func (l *Logger) saveMetrics() {
	// save the metrics to a file
	conf := config.GetInstance()
	file, err := os.OpenFile(conf.MetricsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// the metrics are saved in the following format:
	// timestamp -- log_level -- fastest_chunk -- slowest_chunk -- msg

	_, err = file.WriteString(time.Now().Format(time.RFC3339Nano) + " -- " + l.GetLevel().String() + " -- " + strconv.FormatFloat(l.fastestChunk, 'f', -1, 64) + " -- " + strconv.FormatFloat(l.slowestChunk, 'f', -1, 64) + " -- " + l.metric_name + "\n")
	if err != nil {
		l.Errorf("Error writing to file: %v", err)
	}

	// reset the metrics
	l.fastestChunk = 0
	l.slowestChunk = 0
	l.metric_name = ""

	// close the file
	file.Close()
}

