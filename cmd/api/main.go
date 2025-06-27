package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"risq_backend/config"
	"risq_backend/pkg/app"
	"risq_backend/pkg/logger"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.Log.Level)
	logger.Info("Starting Risk Assessment Backend API...")

	// Create and initialize app
	application := app.New(cfg)
	if err := application.Initialize(); err != nil {
		logger.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in goroutine
	go func() {
		if err := application.Start(); err != nil {
			logger.Errorf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		logger.Info("Received shutdown signal")
	case <-ctx.Done():
		logger.Info("Application context cancelled")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Error during shutdown: %v", err)
	}

	logger.Info("Risk Assessment Backend API stopped")
}
