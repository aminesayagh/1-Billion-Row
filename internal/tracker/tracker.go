package tracker

import (
	"fmt"
	"oneBillion/config"
	"os"
	"runtime"
	"time"
	"sync"
	"sort"
)

type MetricsMap map[string]string

func timer(f func(), metrics chan MetricsMap) {
	// start the timer
	start := time.Now()
	f()

	// end the timer
	end := time.Now()
	diff := end.Sub(start)

	metrics <- MetricsMap{
		"Execution Time": diff.String(),
	}
}

func memory(f func(), metrics chan MetricsMap) {
	// Trigger garbage collection to minimize impact of stale allocations
	runtime.GC()

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

	// Send memory metrics to the metrics channel
	metrics <- map[string]string{
		"Alloc Memory":       fmt.Sprintf("%.2f MB", allocMemory),
		"Total Alloc Memory": fmt.Sprintf("%.2f MB", totalAllocMemory),
		"Sys Memory":         fmt.Sprintf("%.2f MB", sysMemory),
		"Heap Memory":        fmt.Sprintf("%.2f MB", heapMemory),
	}
}

func Run(f func()) {
	done := make(chan bool)
	metrics := make(chan MetricsMap)
	aggregatedMetrics := make(MetricsMap)

	conf := config.GetInstance()

	version := conf.Version
	metricsOutputFilePath := conf.MetricsFilePath

	go func() {
		memory(func() {
			timer(f, metrics)
		}, metrics)

		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("Function execution completed.")
			saveMetrics(metricsOutputFilePath, version, aggregatedMetrics)
			return
		case m := <-metrics:
			for k, v := range m {
				fmt.Printf("--- Metrics %s: %s\n", k, v)
				aggregatedMetrics[k] = v
			}

		case <-time.After(10 * time.Minute):
			fmt.Println("Function execution timed out.")
			saveMetrics(metricsOutputFilePath, version, aggregatedMetrics)
			done <- true
			return
		}
	}
}

func saveMetrics(metricsOutputFilePath string, version string, metrics MetricsMap) {
	file, err := os.OpenFile(metricsOutputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Version: %s\n", version))
	if err != nil {
		fmt.Println(err)
		return
	}

	date := time.Now().Format("2006-01-02 15:04:05")

	for k, v := range metrics {
		_, err = file.WriteString(fmt.Sprintf("%s-%s-%s-%s\n", date, version, k, v))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// MeasureExecutionTime measures the average execution time of two functions and compares them
func MeasureExecutionTime(funcs map[string]func(), iterations int) map[string]time.Duration {
	numFuncs := len(funcs)
	allTimes := make(map[string][]time.Duration, numFuncs)
	allTimesStruct := make(map[string]struct{
		times []time.Duration
		median time.Duration
		p95 time.Duration
		min time.Duration
		max time.Duration
	}, numFuncs)
	var wg sync.WaitGroup
	wg.Add(numFuncs)

	// Measure execution time for each function
	for n, f := range funcs {
		go func(n string, f func()) {
			defer wg.Done()
			var times []time.Duration
			for j := 0; j < iterations; j++ {
				start := time.Now()
				f()
				elapsed := time.Since(start)
				times = append(times, elapsed)
			}
			allTimes[n] = times
		}(n, f)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Process results
	for n, times := range allTimes {
		sort.Slice(times, func(a, b int) bool { return times[a] < times[b] })
		median := times[len(times)/2]
		p95 := times[int(float64(len(times))*0.95)]

		fmt.Printf("Function %s:\n", n)
		fmt.Printf("  Median time: %v\n", median)
		fmt.Printf("  95th percentile time: %v\n", p95)
		fmt.Printf("  Min time: %v\n", times[0])
		fmt.Printf("  Max time: %v\n", times[len(times)-1])

		allTimesStruct[n] = struct{
			times []time.Duration
			median time.Duration
			p95 time.Duration
			min time.Duration
			max time.Duration
		}{
			times: times,
			median: median,
			p95: p95,
			min: times[0],
			max: times[len(times)-1],
		}
	}

	// Return the median times for further analysis
	medianTimes := make(map[string]time.Duration, numFuncs)
	for n, times := range allTimes {
		medianTimes[n] = times[len(times)/2]
	}
	

	return medianTimes
}