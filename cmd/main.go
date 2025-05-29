package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.for.vn/go-web/internal/handlers"
	"go.for.vn/go-web/internal/services"
)

func main() {
	// Create package service
	packageService := services.NewPackageService("data/packages.json")
	if err := packageService.LoadPackages(); err != nil {
		log.Fatal("Failed to load packages:", err)
	}

	// Create handlers
	packageHandler := handlers.NewPackageHandler(packageService)
	apiHandler := handlers.NewAPIHandler(packageService)

	// Setup routes
	http.HandleFunc("/", packageHandler.HandleRoot)
	http.HandleFunc("/api/", apiHandler.HandleAPI)

	// Serve static files directly
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/static"))))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "7666"
	}

	fmt.Printf("🚀 Go Fork Package Registry starting on port %s\n", port)
	fmt.Printf("📦 Visit http://localhost:%s to see package list\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
