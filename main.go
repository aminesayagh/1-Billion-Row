package main

import (
	"bufio"
	"fmt"
	"onBillion/config"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Measurement struct {
	Min    float64
	Max    float64
	Median float64
	Count  int
}

func outputParsing(data map[string]*Measurement) {
	// read the output file
	dataOutputFile, err := os.Create("data/output.csv")
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer dataOutputFile.Close()

	for station, measurement := range data {
		dataOutputFile.WriteString(fmt.Sprintf("%s;%f;%f;%f;%d\n", station, measurement.Min, measurement.Max, measurement.Median, measurement.Count))
	}
}

func hardwareContext() {
	// print some information about the machine running the program
	fmt.Println("--- Hardware Context ---")
	fmt.Println("Number of CPUs: ", runtime.NumCPU())
	fmt.Println("Number of Threads: ", runtime.GOMAXPROCS(0))
	// ram in gegabytes
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated memory: %.2f MB\n", float64(m.Alloc)/1024/1024)            // allocated memory
	fmt.Printf("Total allocated memory: %.2f MB\n", float64(m.TotalAlloc)/1024/1024) // total allocated memory
	fmt.Printf("System memory obtained: %.2f MB\n", float64(m.Sys)/1024/1024)        // system memory obtained
	fmt.Printf("Heap memory in use: %.2f MB\n", float64(m.HeapInuse)/1024/1024)      // heap memory in use, this is memory that has been allocated by the Go runtime to the program
	fmt.Printf("Heap memory released: %.2f MB\n", float64(m.HeapReleased)/1024/1024) // heap memory released, this is memory that has been released by the Go runtime

	fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
	fmt.Println("Number of GoMaxProcs: ", runtime.GOMAXPROCS(0))
}

func softwareContext(file *os.File, config *config.Config) {
	fmt.Println("--- Software Context ---")

	// File system type of the machine
	fileStat, _ := file.Stat()
	var fs syscall.Statfs_t
	syscall.Statfs(fileStat.Name(), &fs)
	fmt.Println("File system type: ", fs.Type)

	// Go version
	fmt.Println("Go version: ", runtime.Version())

	// size of the file
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fmt.Println("File size: ", fileSize)

	// number of lines in the file
	lineCount := config.NumberOfLines
	if lineCount == 0 {
		fileScanner := bufio.NewScanner(file)
		for fileScanner.Scan() {
			lineCount++
		}
	
		// reset the file scanner
		file.Seek(0, 0)
	}
	fmt.Println("Number of lines: ", lineCount)

}

func timer(f func()) {
	// start the timer
	start := time.Now()
	f()
	// end the timer
	end := time.Now()
	diff := end.Sub(start)

	fmt.Println("Execution time: " + diff.String())
}
func memory(f func()) {
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

	fmt.Printf("Memory allocated during function execution: %.2f MB\n", allocMemory)
	fmt.Printf("Total memory allocated by the function: %.2f MB\n", totalAllocMemory)
	fmt.Printf("System memory used by the function: %.2f MB\n", sysMemory)
	fmt.Printf("Heap memory used by the function: %.2f MB\n", heapMemory)
}
func tracker(f func()) {
	done := make(chan bool)

	go func() {
		memory(func() {
			timer(f)
			done <- true
		})
	}()

	select {
	case <-done:
		fmt.Println("Function execution completed.")
	case <-time.After(10 * time.Minute):
		fmt.Println("Function execution timed out.")
	}
}

func parsing() {
	// read the input env
	config := config.GetInstance()

	// read the input file
	dataFile, err := os.Open(config.InputFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer dataFile.Close()

	// file scanning logic
	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	// parsing logic
	measurements := make(map[string]*Measurement)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			continue
		}
		station := parts[0]

		temperature, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			fmt.Println("Error parsing temperature: ", err)
			continue
		}

		measurement, exists := measurements[station]

		if !exists {
			measurement = &Measurement{
				Min:    temperature,
				Max:    temperature,
				Median: temperature,
				Count:  1,
			}
			measurements[station] = measurement
		} else {
			if temperature < measurement.Min {
				measurement.Min = temperature
			} else if temperature > measurement.Max {
				measurement.Max = temperature
			}
			measurement.Median = (measurement.Median + temperature) / 2
			measurement.Count++
		}
	}

	outputParsing(measurements)
}

func main() {
	// High performance data processing pipeline
	config := config.GetInstance()

	fmt.Println("Input file path: ", config.InputFilePath)

	dataFile, err := os.Open(config.InputFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer dataFile.Close()

	// file scanning logic
	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	// hardware context
	hardwareContext()

	// software context
	softwareContext(dataFile, config)

	// close the file
	dataFile.Close()
	// remove the variable
	dataFile = nil

	tracker(parsing)
}
