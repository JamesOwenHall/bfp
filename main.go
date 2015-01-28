package main

import (
	"fmt"
	"github.com/JamesOwenHall/BruteForceProtection/server"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	// Setup multithreading
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create server
	s := server.New()
	defer s.Close()

	// Setup routes
	s.Routes["password"] = func(value interface{}) server.Response {
		return server.Response{
			Valid: true,
		}
	}

	// Start server
	go s.ListenAndServe("/tmp/bfp.sock")
	fmt.Println("Now listening at /tmp/bfp.sock")

	// Capture interrupt signal so that the server closes properly
	interrupts := make(chan os.Signal)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts
	fmt.Println("Shutting down")
}
