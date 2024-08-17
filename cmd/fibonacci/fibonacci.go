package main


import (
	"fmt"
	"onBillion/internal/logger"
)

func fibonacci(n int, memo map[int]int) int {
	_, ok := memo[n]
	if ok {
		return memo[n]
	}
	if n <= 2 {
		return 1
	}
	memo[n] = fibonacci(n-1, memo) + fibonacci(n-2, memo)

	return memo[n]
}

func main() {
	log := logger.NewLogger("info")
	log.TrackPerformance(func() {
		n := 1000000
		memo := make(map[int]int)

		result := fibonacci(n, memo)
		fmt.Println(result)
	}, "Fibonacci Performance")
}

