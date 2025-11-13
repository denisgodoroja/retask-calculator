package main

import (
	"log"
	"net/http"
	"os"

	"denisgodoroja/retask/internal/service"
	"denisgodoroja/retask/internal/storage/inmemory"
	"denisgodoroja/retask/internal/webservice"
)

func main() {
	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local development
	}

	// Create the In-Memory repository
	repo := inmemory.NewInMemoryPackRepo()

	// Create the service layer
	packService := service.NewPackService(repo)

	// Create the HTTP handler layer
	handler := webservice.NewHandler(packService)

	// Create the router
	router := webservice.NewRouter(handler)

	log.Printf("Starting Go backend server on http://localhost:%s", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start listening for incoming HTTP requests
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
