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
	configFilename := flag.String("c", "config.json", "the name of the configuration file")
	flag.Parse()

	// Setup multithreading
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Read the configuration
	configuration, errs := config.ReadConfig(*configFilename)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, "configuration error:", err)
		}

		return
	}

	// Start server
	counter := hitcounter.NewHitCounter(configuration.Directions)
	defer counter.Close()
	go counter.ListenAndServe(configuration.ListenType, configuration.ListenAddress)
	fmt.Println("Now listening at", configuration.ListenAddress)

	// Capture interrupt signal so that the server closes properly
	interrupts := make(chan os.Signal)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts
}
