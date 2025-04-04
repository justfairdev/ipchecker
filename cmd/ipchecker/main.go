package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/justfairdev/ipchecker/internal/config"
	"github.com/justfairdev/ipchecker/internal/server"
)

// main is the entry point for the IPChecker application.
//
// Application Overview:
//   - Loads configuration settings (ports, database paths, etc.).
//   - Initializes combined HTTP (Gin) and gRPC servers along with shared dependencies.
//   - Starts the servers concurrently, making services available to HTTP and gRPC clients.
//   - Gracefully handles system interrupts (SIGINT, SIGTERM) to safely shut down servers.
//
// This structure allows the application to serve multiple client types concurrently, manage graceful shutdown,
// and provides clear logging for observability and debugging.
func main() {
	// Load application configuration from environment variables, files, or defaults.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize HTTP/gRPC AppServer with shared dependencies (e.g., GeoLookup database).
	appServer, err := server.NewAppServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create AppServer: %v", err)
	}

	// Start the combined HTTP and gRPC servers concurrently in a separate goroutine.
	// HTTP listens on the port defined in cfg.HTTPPort, gRPC listens on port "50051" (can be customized).
	go func() {
		if err := appServer.Start(cfg.HTTPPort, "50051"); err != nil {
			log.Fatalf("Server encountered an error during startup: %v", err)
		}
	}()

	// Set up OS signal channel to listen for termination signals (Ctrl+C, Docker/Kubernetes shutdown, etc.).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block execution until a shutdown signal is received.
	<-quit

	log.Println("Shutdown signal received, gracefully stopping servers...")

	// Gracefully shut down both gRPC server and safely release resources (GeoLookup database connections, etc.).
	appServer.Stop()

	log.Println("All servers stopped successfully. Exiting.")
}
