package main

import (
	"bufio"
	"oneBillion/cmd/version/v1_base"
	"fmt"
	"oneBillion/config"
	"oneBillion/internal/context"
	"oneBillion/internal/tracker"
	"os"
)


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
		v1base.Parsing(config)
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
