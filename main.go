package main

import (
	"fmt"
	"github.com/JamesOwenHall/BruteForceProtection/hitcounter"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	// Setup multithreading
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setup directions
	directions := []hitcounter.Direction{
		hitcounter.NewStringDirection("password", 10, 3),
		hitcounter.NewInt32Direction("id", 10, 3),
	}

	// Start server
	counter := hitcounter.NewHitCounter(directions)
	defer counter.Close()
	go counter.ListenAndServe("/tmp/bfp.sock")
	fmt.Println("Now listening at /tmp/bfp.sock")

	// Capture interrupt signal so that the server closes properly
	interrupts := make(chan os.Signal)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts
	fmt.Println("Shutting down")
}
