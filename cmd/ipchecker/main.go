package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/justfairdev/ipchecker/internal/config"
	"github.com/justfairdev/ipchecker/internal/server"
)

func main() {
	// 1) Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2) Create the combined server
	appServer, err := server.NewAppServer(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// 3) Start servers in a separate goroutine so we can watch for shutdown signals
	go func() {
		// HTTP on cfg.HTTPPort, gRPC on 50051 (hard-coded or from config).
		if err := appServer.Start(cfg.HTTPPort, "50051"); err != nil {
			log.Fatalf("failed to start servers: %v", err)
		}
	}()

	// 4) Listen for OS signals (Ctrl+C, Docker stop, etc.) for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	appServer.Stop()
	log.Println("Done.")
}
