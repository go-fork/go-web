package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-fork/go-web/internal/handlers"
	"github.com/go-fork/go-web/internal/services"
)

func main() {
	// Initialize package service
	packageService := services.NewPackageService("data/packages.json")

	// Initialize handlers
	handler := handlers.NewHandler(packageService)

	// Static files - only if public directory exists
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Routes
	// Main handler that routes all requests
	http.HandleFunc("/", handler.RouteRequest)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
