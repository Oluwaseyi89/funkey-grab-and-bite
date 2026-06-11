//go:build !lambda

package main

import (
	"funkey-grab-and-bite/funkey-bite-api/internal/app"
	"log"
)

func main() {
	log.Println("Starting Funkey Grab-and-Bite API in Standard LOCAL HTTP Mode...")

	// Initialize the shared engine configuration
	router, cleanup := app.SetupEngine()
	defer cleanup()

	// Launch engine on the standard server TCP port listener
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Local server encountered unrecoverable runtime crash: %v", err)
	}
}