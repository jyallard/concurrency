package main

import (
	"fmt"
	"sync"
	"time"
)

var iterations = 500000
var sharedCounter = 0
var mutex = &sync.Mutex{}

func updateCounter(del int) {
	mutex.Lock()
	sharedCounter += del // Low-level ops are made atomic
	mutex.Unlock()
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

	fmt.Println("Mutex-based sync final count should be 0: ", sharedCounter)
	fmt.Println("Elapsed: ", duration)
}
