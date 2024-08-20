package tracker

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"onBillion/config"
)

type MetricsMap map[string]string

func trackMetrics(ctx context.Context, f func(), metrics chan MetricsMap) {
	start := time.Now()

	// Memory statistics before executing the function
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	f()

	// Memory statistics after executing the function
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calculate memory used by the function
	allocMemory := float64(memAfter.Alloc-memBefore.Alloc) / 1024 / 1024
	totalAllocMemory := float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / 1024 / 1024
	sysMemory := float64(memAfter.Sys-memBefore.Sys) / 1024 / 1024
	heapMemory := float64(memAfter.HeapAlloc-memBefore.HeapAlloc) / 1024 / 1024

	// End the timer
	end := time.Now()
	diff := end.Sub(start)

	// Send all metrics to the metrics channel
	metrics <- map[string]string{
		"ExecutionTime":    diff.String(),
		"AllocMemory":      fmt.Sprintf("%.2f MB", allocMemory),
		"TotalAllocMemory": fmt.Sprintf("%.2f MB", totalAllocMemory),
		"SysMemory":        fmt.Sprintf("%.2f MB", sysMemory),
		"HeapMemory":       fmt.Sprintf("%.2f MB", heapMemory),
		// add more metrics here
	}

	// Optionally, you can monitor the context for cancellation or timeout
	select {
	case <-ctx.Done():
		fmt.Println("Function execution was cancelled or timed out.")
	default:
		// continue normally
	}
}

func Run(f func()) {
	var wg sync.WaitGroup
	metrics := make(chan MetricsMap)
	aggregatedMetrics := make(MetricsMap)

	conf := config.GetInstance()
	version := conf.Version
	metricsOutputFilePath := conf.MetricsFilePath

	// Use a context with timeout to manage function execution time
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		trackMetrics(ctx, f, metrics)
	}()

	go func() {
		for m := range metrics {
			for k, v := range m {
				fmt.Printf("%s: %s\n", k, v)
				aggregatedMetrics[k] = v
			}
		}
	}()

	wg.Wait()
	close(metrics)
	saveMetrics(metricsOutputFilePath, version, aggregatedMetrics)
	fmt.Println("Function execution completed.")
}

func saveMetrics(metricsOutputFilePath string, version string, metrics MetricsMap) {
	tempFilePath := metricsOutputFilePath + ".tmp"

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("Error creating temporary metrics file: ", err)
		return
	}
	defer tempFile.Close()

	// save the metrics to the file in the format: date-version-key=value
	date := time.Now().Format("2006-01-02")
	for key, value := range metrics {
		if hasMetric, _ := searchMetric(metricsOutputFilePath, version, key); hasMetric != 0 {
			// remove the line from original file
			removeMetric(metricsOutputFilePath, hasMetric)
		}
		_, _ = tempFile.WriteString(fmt.Sprintf("%s-%s-%s=%s\n", date, version, key, value))
	}

	// Rename temp file to actual file
	if err := os.Rename(tempFilePath, metricsOutputFilePath); err != nil {
		fmt.Println("Error renaming temporary metrics file: ", err)
	} else {
		fmt.Println("Metrics saved to file: ", metricsOutputFilePath)
	}
}

func searchMetric(metricsOutputFilePath, version, key string) (int, string) {
	metricsOutputFile, err := os.Open(metricsOutputFilePath)
	if err != nil {
		fmt.Println("Error opening metrics file: ", err)
		return 0, ""
	}
	defer metricsOutputFile.Close()

	scanner := bufio.NewScanner(metricsOutputFile)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		line := scanner.Text()
		parts := strings.Split(line, "-")
		if len(parts) < 3 {
			continue
		}
		if parts[1] == version && strings.Contains(parts[2], key) {
			return currentLine, line
		}
	}
	return 0, ""
}

func removeMetric(metricsOutputFilePath string, line int) {
	metricsOutputFile, err := os.Open(metricsOutputFilePath)
	if err != nil {
		fmt.Println("Error opening metrics file: ", err)
		return
	}
	defer metricsOutputFile.Close()

	tempFilePath := metricsOutputFilePath + ".tmp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("Error creating temporary file: ", err)
		return
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(metricsOutputFile)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine == line {
			continue
		}
		_, _ = tempFile.WriteString(scanner.Text() + "\n")
	}

	if err := os.Rename(tempFilePath, metricsOutputFilePath); err != nil {
		fmt.Println("Error renaming temporary file: ", err)
	}
}
