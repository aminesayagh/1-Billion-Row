package version_2_concurrency

import (
	"bufio"
	"fmt"
	"onBillion/config"
	"os"
	"time"
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

	// reset the file scanner
	dataFile.Seek(0, 0)

	fileScanner = bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	// start a timer to measure the time taken to process the file
	start := time.Now()

	chunkNumber := config.ChunkSize
	chunk := make(chan [string]*Measurement, chunkNumber)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		chunk = append(chunk, line)
		if len(chunk) == chunkNumber {
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
