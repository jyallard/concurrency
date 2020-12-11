package main

import (
	"fmt"
	"sync"
	"time"
)

var iterations = 500000
var sharedCounter = 0

func updateCounter(del int) {
	sharedCounter += del // A the low level, this is done by at least three ops
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			updateCounter(1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			updateCounter(-1)
		}
	}()

	wg.Wait() // Wait for all routines to finish
	duration := time.Since(start)

	fmt.Println("Unpredictable final count should be 0: ", sharedCounter)
	fmt.Println("Elapsed: ", duration)
}
