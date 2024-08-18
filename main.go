package main


import (
	"bufio"
	"onBillion/config"
	"time"
	"syscall"
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

	
	// file scanning logic
	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)
	
	// hardware context
	hardwareContext()
	
	// software context
	softwareContext(dataFile)
	
	fileScanner = bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)
	
	timer(func() {
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
				}
				if temperature > measurement.Max {
					measurement.Max = temperature
				}
				measurement.Median = (measurement.Median + temperature) / 2
				measurement.Count++
			}
		}

		outputParsing(measurements)
	})

}

func hardwareContext() {
	// print some information about the machine running the program
	fmt.Println("--- Hardware Context ---")
	fmt.Println("Number of CPUs: ", runtime.NumCPU())
	fmt.Println("Number of Threads: ", runtime.GOMAXPROCS(0))
	// ram in gegabytes
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("Allocated memory: ", m.Alloc/1024/1024, "MB")
	fmt.Println("Total memory: ", m.TotalAlloc/1024/1024, "MB")
	
	fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
	fmt.Println("Number of GoMaxProcs: ", runtime.GOMAXPROCS(0))
}

func softwareContext(file *os.File) {
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
	lineCount := 0
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		lineCount++
	}
	fmt.Println("Number of lines: ", lineCount) // number of lines in the file is 0, because the file has been read to the end of the file

	// reset the file scanner
	file.Seek(0, 0)
}

func timer(f func()) {
	done := make(chan bool)

	go func() {
		// start the timer
		start := time.Now()
		f()
		// end the timer
		end := time.Now()
		diff := end.Sub(start)

		fmt.Println("Execution time: " + diff.String())
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("Function execution completed.")
	case <-time.After(10 * time.Minute):
		fmt.Println("Function execution timed out.")
	}
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