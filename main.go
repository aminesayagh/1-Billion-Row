package main

import (
	"fmt"
	"onBillion/config"
	"onBillion/internal/logger"
)



func main() {
	fmt.Println("Hello, World!")

	config := config.GetInstance()

	fmt.Println(config.InputFilePath)
	
	logger := logger.NewLogger("info");
	logger.TrackPerformance(func() {
		
	}, "Main Performance");
}
