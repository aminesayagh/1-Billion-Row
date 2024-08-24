package main

import (
	"bufio"
	"bytes"
	"fmt"
	"onBillion/config"
	"onBillion/internal/byteutil"
	"onBillion/internal/context"
	"onBillion/internal/map"
	"onBillion/internal/tracker"
	"os"
)

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

	// map to store the measurements
	measurements := mapStruct.New()

	var (
		currentSepIndex int
	)

	for fileScanner.Scan() {
		currentLine := fileScanner.Bytes()

		currentSepIndex = bytes.IndexByte(currentLine, ';') // split by semicolon without memory allocation

		if currentSepIndex == -1 {
			continue
		}

		temperature := byteutil.ConvertTo4Byte(currentLine[currentSepIndex+1:])
		measurements.Add(currentLine[:currentSepIndex], temperature)
	}

	dataOutputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer dataOutputFile.Close()

	// Write output efficiently
	outputBuffer := bufio.NewWriter(dataOutputFile)
	for _, measurement := range measurements.List() {
		outputBuffer.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s\n", measurement.Name, measurement.Min, measurement.Max, measurement.Average, measurement.Count))
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

	tracker.Run(parsing)

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
