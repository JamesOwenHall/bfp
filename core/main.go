package main

import (
	"flag"
	"fmt"
	"github.com/JamesOwenHall/BruteForceProtection/core/config"
	"github.com/JamesOwenHall/BruteForceProtection/core/hitcounter"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	// Parse flags
	configFilename := flag.String("c", "bfp-config.json", "the name of the configuration file")
	flag.Parse()

	// Setup multithreading
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Read the configuration
	configuration, err := config.ReadConfig(*configFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}

	// Start server
	counter := hitcounter.NewHitCounter(configuration.Directions)
	defer counter.Close()
	go counter.ListenAndServe(configuration.ListenAddress)
	fmt.Println("Now listening at", configuration.ListenAddress)

	// Capture interrupt signal so that the server closes properly
	interrupts := make(chan os.Signal)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts
}
