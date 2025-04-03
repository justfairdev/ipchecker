package main

import (
    "log"

    "github.com/justfairdev/ipchecker/internal/config"
    "github.com/justfairdev/ipchecker/internal/server"
)

func main() {
    // 1. Load config (environment vars or defaults)
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // 2. Create the server
    srv, err := server.NewServer(cfg)
    if err != nil {
        log.Fatalf("failed to create server: %v", err)
    }

    // 3. Run the server on the configured port
    if err := srv.Run(":" + cfg.HTTPPort); err != nil {
        log.Fatalf("failed to run server: %v", err)
    }
}
