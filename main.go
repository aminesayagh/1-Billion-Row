package main


import (
	"bufio"
	"onBillion/config"
	"time"
	"strings"
	"os"
	"strconv"
	"fmt"
	"runtime"
)

type Measurement struct {
	Min    float64
	Max    float64
	Median float64
	Count  int
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

	measurements := make(map[string]*Measurement)
	
	// file scanning logic
	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	// size of the file
	fileInfo, _ := dataFile.Stat()
	fileSize := fileInfo.Size()
	fmt.Println("File size: ", fileSize)

	// number of lines in the file
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	fmt.Println("Number of lines: ", lineCount) // number of lines in the file is 0, because the file has been read to the end of the file

	// print some information about the machine running the program
	fmt.Println("Number of CPUs: ", runtime.NumCPU())
	fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
	fmt.Println("Number of GoMaxProcs: ", runtime.GOMAXPROCS(0))
	

	// reset the file scanner
	dataFile.Seek(0, 0)

	fileScanner = bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	// start a timer to measure the time taken to process the file
	start := time.Now()
	// read the file line by line
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
			}
			if temperature > measurement.Max {
				measurement.Max = temperature
			}
			measurement.Median = (measurement.Median + temperature) / 2
			measurement.Count++
		}
	}
	// stop the timer
	elapsed := time.Since(start)
	fmt.Println("Time taken: ", elapsed)

	outputParsing(measurements)
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