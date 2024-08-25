package main

import (
	"bufio"
	"bytes"
	"fmt"
	"onBillion/config"
	"onBillion/internal/context"
	"onBillion/internal/tracker"
	"os"
	"strconv"
	"strings"
)

type Measurement struct {
	Min     int16
	Max     int16
	Average int16
	Count   int32
}

func ConvertByteToInt(b []byte, cache map[string]int16) (int16, error) {
	str := strings.TrimRight(string(b), ".")
	str = strings.SplitN(str, ".", 2)[0]

	if cachedValue, found := cache[str]; found {
		return cachedValue, nil
	}

	// Convert string to integer using strconv.Atoi (as we're handling integers)
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("error converting %s to int: %v", str, err)
	}

	if i < -32768 || i > 32767 {
		return 0, fmt.Errorf("integer out of int16 range: %d", i)
	}

	int16Value := int16(i)
	cache[str] = int16Value
	return int16Value, nil
}

func parsing(config *config.Config) {
	// read the input file
	dataFile, err := os.Open(config.InputFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer dataFile.Close()

	// file scanning logic
	fileScanner := bufio.NewScanner(dataFile)
	buf := make([]byte, 0, 256*1024)   // 64 KB buffer
	fileScanner.Buffer(buf, 1024*1024) // Up to 1 MB lines

	measurements := make(map[string]*Measurement, 500000)
	temperatureCache := make(map[string]int16, 400)

	for fileScanner.Scan() {
		currentLine := fileScanner.Bytes()

		currentSepIndex := bytes.IndexByte(currentLine, ';') // split by semicolon without memory allocation

		if currentSepIndex == -1 {
			continue
		}
		currentStation := string(currentLine[:currentSepIndex])
		currentTemperature, err := ConvertByteToInt(currentLine[currentSepIndex+1:], temperatureCache)
		if err != nil {
			fmt.Println("Error converting temperature: ", err)
			continue
		}

		if currentMeasurement, exists := measurements[currentStation]; !exists {
			measurements[currentStation] = &Measurement{
				Min:     currentTemperature,
				Max:     currentTemperature,
				Average: currentTemperature,
				Count:   1,
			}
		} else {
			if currentTemperature > currentMeasurement.Max {
				currentMeasurement.Max = currentTemperature
			}
			if currentTemperature < currentMeasurement.Min {
				currentMeasurement.Min = currentTemperature
			}
			currentMeasurement.Count++
			currentMeasurement.Average = int16((int32(currentMeasurement.Average)*(currentMeasurement.Count-1) + int32(currentTemperature)) / currentMeasurement.Count)
		}
	}

	dataOutputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer dataOutputFile.Close()

	// Write output efficiently
	outputBuffer := bufio.NewWriterSize(dataOutputFile, 64*1024) // 64 KB buffer

	for station, measurement := range measurements {
		fmt.Fprintf(outputBuffer, "%s;%d;%d;%d;%d\n", station, measurement.Min, measurement.Max, measurement.Average, measurement.Count)
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

	tracker.Run(func() {
		parsing(config)
	})

	// print the output file 10 lines for testing
	dataOutputFile, err := os.Open(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}

	defer dataOutputFile.Close()

	fileScanner = bufio.NewScanner(dataOutputFile)
	fileScanner.Split(bufio.ScanLines)

	for i := 0; i < 10; i++ {
		fileScanner.Scan()
		fmt.Println(fileScanner.Text())
	}
}
