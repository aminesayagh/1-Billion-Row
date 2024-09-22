package v2_assembly

import (
	"bufio"
	"bytes"
	"fmt"
	"oneBillion/config"
	"os"
	"strconv"
)

type Measurement struct {
	Min     DecimalNumber
	Max     DecimalNumber
	Average DecimalNumber
    Sum     DecimalNumber
	Count   int32
}

// Assembly function declaration
// BytesToNumericBytes parses the input byte slice and writes valid numeric bytes to the output slice.
// Returns the number of bytes written or a negative error code.
func BytesToNumericBytes(input []byte, output []byte) int

func Parsing(config *config.Config) {
	// Open the input file
	dataFile, err := os.Open(config.InputFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer dataFile.Close()

	// Initialize the file scanner
	fileScanner := bufio.NewScanner(dataFile)
	buf := make([]byte, 0, 64*1024)    // 64 KB buffer
	fileScanner.Buffer(buf, 1024*1024) // Up to 1 MB lines

	// Initialize data structures
	measurements := make(map[string]*Measurement, 500000)
	temperatureCache := make(map[string]DecimalNumber, 10000)

	var (
		currentLine     []byte
		currentSepIndex int
		currentStation  string
		tempBytes       []byte
		numericBytes    []byte
		tempKey         string
		currentTemp     DecimalNumber
		exists          bool
	)

	// Pre-allocate buffers for processing temperatures
	maxLineLength := 12 // Adjust based on expected maximum line length
	tempBytes = make([]byte, maxLineLength)
	numericBytes = make([]byte, maxLineLength)

	for fileScanner.Scan() {
		currentLine = fileScanner.Bytes()

		currentSepIndex = bytes.IndexByte(currentLine, ';') // Split by semicolon without memory allocation
		if currentSepIndex == -1 {
			continue
		}

		// Extract the station as a string (key for the map)
		currentStation = string(currentLine[:currentSepIndex])

		// Extract the temperature bytes
		tempBytes = currentLine[currentSepIndex+1:]

		// Use the temperature bytes as a key for caching
		tempKey = string(tempBytes)

		if currentTemp, exists = temperatureCache[tempKey]; !exists {
            l := len(tempBytes)
            
			n := BytesToNumericBytes(tempBytes, numericBytes)
			if n < 0 {
				fmt.Println("Error parsing temperature: ", n)
				continue
			}

            // redimension numericBytes by the value of l

			// Convert the numeric bytes to a float64
            var newTemp DecimalNumber
            newTemp.Normalize(numericBytes)
            currentTemp = newTemp
			temperatureCache[tempKey] = currentTemp
		}

		// Update measurements
		if currentMeasurement, exists := measurements[currentStation]; exists {
            
			if currentTemp.Compare(currentMeasurement.Max) > 0 {
				currentMeasurement.Max = currentTemp
			}
			if currentTemp.Compare(currentMeasurement.Min) < 0 {
				currentMeasurement.Min = currentTemp
			}
			currentMeasurement.Count++
            currentMeasurement.Sum = currentMeasurement.Sum.Add(currentTemp)
		} else {
			measurements[currentStation] = &Measurement{
				Min:     currentTemp,
				Max:     currentTemp,
				Average: currentTemp,
                Sum:     currentTemp,
				Count:   1,
			}
		}
	}

	if err := fileScanner.Err(); err != nil {
		fmt.Println("Error scanning file: ", err)
	}

	// Write output to file
	outputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer outputFile.Close()

	outputBuffer := bufio.NewWriterSize(outputFile, 4*1024*1024) // 4 MB buffer
	lineBuffer := make([]byte, 0, 128)

    var average DecimalNumber

	for station, measurement := range measurements {
		lineBuffer = lineBuffer[:0] // Reset line buffer
        average = CalculateAverage(measurement.Sum, measurement.Count)

		lineBuffer = append(lineBuffer, station...)
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
