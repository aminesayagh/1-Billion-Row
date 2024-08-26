package v1base

import (
	"bufio"
	"fmt"
	"oneBillion/config"
	"os"
	"strconv"
	"bytes"
)

type Measurement struct {
	Min     float64
	Max     float64
	Average float64
	Count   int32
}

func Parsing(config *config.Config) {
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

	measurements := make(map[string]*Measurement, 500000)
	temperatureCache := make(map[string]float64, 10000)

	var (
		currentLine              []byte
		currentSepIndex          int
		currentStation           string
		currentTemperatureString string
		currentTemperature       float64
		exists                   bool
	)

	
	for fileScanner.Scan() {
		currentLine = fileScanner.Bytes()

		currentSepIndex = bytes.IndexByte(currentLine, ';') // split by semicolon without memory allocation
		if currentSepIndex == -1 {
			continue
		}

		currentStation = string(currentLine[:currentSepIndex])
		currentTemperatureString = string(currentLine[currentSepIndex+1:])

		if currentTemperature, exists = temperatureCache[currentTemperatureString]; !exists {
			currentTemperature, err = strconv.ParseFloat(currentTemperatureString, 64)
			if err != nil {
				fmt.Println("Error parsing temperature: ", err)
				continue
			}
			temperatureCache[currentTemperatureString] = currentTemperature
		}

		if currentMeasurement, exists := measurements[currentStation]; exists {
			if currentTemperature > currentMeasurement.Max {
				currentMeasurement.Max = currentTemperature
			}
			if currentTemperature < currentMeasurement.Min {
				currentMeasurement.Min = currentTemperature
			}
			currentMeasurement.Count++
			currentMeasurement.Average = ((currentMeasurement.Average * float64(currentMeasurement.Count - 1)) + currentTemperature) / float64(currentMeasurement.Count)
		} else {
			measurements[currentStation] = &Measurement{
				Min:     currentTemperature,
				Max:     currentTemperature,
				Average: currentTemperature,
				Count:   1,
			}
		}
	}
	outputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer outputFile.Close()

	// Write output efficiently
	outputBuffer := bufio.NewWriterSize(outputFile, 4*1024*1024) // 64 KB buffer
	lineBuffer := make([]byte, 0, 64)

	for station, measurement := range measurements {
		lineBuffer = append(lineBuffer[:0], station...)
		lineBuffer = append(lineBuffer, ';')
		lineBuffer = strconv.AppendFloat(lineBuffer, measurement.Min, 'f', 2, 64)
		lineBuffer = append(lineBuffer, ';')
		lineBuffer = strconv.AppendFloat(lineBuffer, measurement.Max, 'f', 2, 64)
		lineBuffer = append(lineBuffer, ';')
		lineBuffer = strconv.AppendFloat(lineBuffer, measurement.Average, 'f', 2, 64)
		lineBuffer = append(lineBuffer, ';')
		lineBuffer = strconv.AppendInt(lineBuffer, int64(measurement.Count), 10)
		lineBuffer = append(lineBuffer, '\n')

		_, err := outputBuffer.Write(lineBuffer)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
	if err := outputBuffer.Flush(); err != nil {
		fmt.Println("Error flushing output buffer:", err)
	}
}