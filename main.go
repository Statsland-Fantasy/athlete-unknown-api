//go:build !lambda
// +build !lambda

package main

import (
	"log"
	"os"
)

func main() {
	router := SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API endpoints available at /v1/*")
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
