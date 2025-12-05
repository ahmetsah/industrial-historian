package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/config"
	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/core"
	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/repository"
	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/transport"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	if err := run(ctx); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("Starting Alarm Service on port %s", cfg.Port)

	// Initialize Postgres Repository
	repo, err := repository.NewPostgresRepository(cfg.DbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer repo.Close()

	// Initialize NATS Transport
	natsTransport, err := transport.NewNatsTransport(cfg.NatsUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer natsTransport.Close()

	// Initialize Alarm Service
	// NatsTransport implements EventPublisher
	svc := core.NewAlarmService(repo, natsTransport)

	// Set service in NatsTransport (for consumer)
	natsTransport.SetService(svc)

	// Load definitions from DB
	if err := svc.LoadDefinitions(); err != nil {
		return fmt.Errorf("failed to load definitions: %w", err)
	}

	// Start NATS Consumer
	if err := natsTransport.Start(); err != nil {
		return fmt.Errorf("failed to start NATS consumer: %w", err)
	}

	// Initialize HTTP Handler
	httpHandler := transport.NewHttpHandler(svc)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Start HTTP Server
	go func() {
		log.Printf("HTTP server listening on %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Shutting down Alarm Service...")

	// Shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown failed: %v", err)
	}

	return nil
}
