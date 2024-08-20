package context

import (
	"bufio"
	"fmt"
	"onBillion/config"
	"os"
	"runtime"
	"syscall"
)

func HardwareContext() {
	// print some information about the machine running the program
	fmt.Println("--- Hardware Context ---")
	fmt.Println("Number of CPUs: ", runtime.NumCPU())
	fmt.Println("Number of Threads: ", runtime.GOMAXPROCS(0))
	// ram in gigabytes
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

func SoftwareContext(file *os.File, config *config.Config) {
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
	lineCount := 1000010000
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
