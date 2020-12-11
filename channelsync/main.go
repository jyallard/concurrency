package main

import (
	"fmt"
	"sync"
	"time"
)

type counterOp struct {
	delta   int
	confirm chan int // confirmation channel
}

var iterations = 500000
var sharedCounter = 0
var counterRequests chan *counterOp

func updateCounter(del int) int {
	update := &counterOp{delta: del, confirm: make(chan int)}
	counterRequests <- update
	newSharedCounter := <-update.confirm

	return newSharedCounter
}

func main() {
	counterRequests = make(chan *counterOp, 10)
	var wg sync.WaitGroup
	start := time.Now()

	go func() {
		for {
			select {
			case request := <-counterRequests:
				sharedCounter += request.delta
				request.confirm <- sharedCounter
			}
		}
	}()

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

	fmt.Println("Channel-based sync final count should be 0: ", sharedCounter)
	fmt.Println("Elapsed: ", duration)
}
