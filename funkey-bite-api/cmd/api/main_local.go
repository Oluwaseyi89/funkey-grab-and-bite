//go:build !lambda

package main

import (
	"funkey-grab-and-bite/funkey-bite-api/internal/app"
	"log"
	"os"
)

func main() {
	log.Println("Starting Funkey Grab-and-Bite API in Standard LOCAL HTTP Mode...")

	// Initialize the shared engine configuration
	router, cleanup := app.SetupEngine()
	defer cleanup()

	// Host platforms like Railway assign the listen port dynamically via $PORT.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Launch engine on the standard server TCP port listener
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Local server encountered unrecoverable runtime crash: %v", err)
	}
}