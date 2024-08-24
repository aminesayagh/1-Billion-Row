package main

import (
	"bufio"
	"fmt"
	"onBillion/config"
	"onBillion/internal/context"
	"onBillion/internal/tracker"
	"os"
	"strconv"
	"strings"
)

type Measurement struct {
	Min     float64
	Max     float64
	Average float64
	Count   int
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
	buf := make([]byte, 0, 64*1024)    // 64 KB buffer
	fileScanner.Buffer(buf, 1024*1024) // Up to 1 MB lines

	// parsing logic
	measurements := make(map[string]*Measurement, 100000)

	currentLine := ""
	currentStation := ""
	currentTemperature := ""
	currentSepIndex := -1
	currentMeasurement := &Measurement{}
	currentMeasurementExists := false

	for fileScanner.Scan() {
		currentLine = fileScanner.Text()

		currentSepIndex = strings.Index(currentLine, ";") // split by semicolon without memory allocation

		if currentSepIndex == -1 {
			continue
		}

		currentStation = currentLine[:currentSepIndex]
		currentTemperature = currentLine[currentSepIndex+1:]

		parsedTemperature, err := strconv.ParseFloat(currentTemperature, 64)
		if err != nil {
			fmt.Println("Error parsing temperature: ", err)
			continue
		}

		currentMeasurement, currentMeasurementExists = measurements[currentStation]

		if !currentMeasurementExists {
			currentMeasurement = &Measurement{
				Min:     parsedTemperature,
				Max:     parsedTemperature,
				Average: parsedTemperature,
				Count:   1,
			}
			measurements[currentStation] = currentMeasurement
		} else {
			if parsedTemperature < currentMeasurement.Min {
				currentMeasurement.Min = parsedTemperature
			}
			if parsedTemperature > currentMeasurement.Max {
				currentMeasurement.Max = parsedTemperature
			}
			currentMeasurement.Count++
			currentMeasurement.Average = (parsedTemperature - currentMeasurement.Average) / float64(currentMeasurement.Count)
		}
	}

	dataOutputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer dataOutputFile.Close()

	// Write output efficiently
	outputBuffer := bufio.NewWriter(dataOutputFile)
	for station, measurement := range measurements {
		outputBuffer.WriteString(station + ";" +
			strconv.FormatFloat(measurement.Min, 'f', 6, 64) + ";" +
			strconv.FormatFloat(measurement.Max, 'f', 6, 64) + ";" +
			strconv.FormatFloat(measurement.Average, 'f', 6, 64) + ";" +
			strconv.Itoa(measurement.Count) + "\n")
	}
	outputBuffer.Flush()
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
	context.HardwareContext()

	// software context
	context.SoftwareContext(dataFile, config)

	// close the file
	dataFile.Close()
	// remove the variable
	dataFile = nil

	// parsing
	tracker.Run(parsing)
}
