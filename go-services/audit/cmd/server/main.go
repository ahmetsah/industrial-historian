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

	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/config"
	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/core"
	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/repository"
	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/transport"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("Starting Audit Service with DB: %s, NATS: %s", cfg.DbUrl, cfg.NatsUrl)

	// Context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, err := repository.NewPostgresRepository(ctx, cfg.DbUrl)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	defer repo.Close()

	hasher := core.NewSHA256Hasher()

	consumer, err := transport.NewAuditConsumer(cfg.NatsUrl, repo, hasher)
	if err != nil {
		return fmt.Errorf("failed to initialize consumer: %w", err)
	}
	defer consumer.Close()

	if err := consumer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	// HTTP Server
	httpHandler := transport.NewHttpHandler(repo, hasher)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/audit/verify", httpHandler.Verify)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting HTTP server on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server failed: %v", err)
		}
	}()

	log.Println("Audit Service running...")

	// Wait for signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	// Cleanup logic here (close DB, NATS, etc.)
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	// Defer calls handle closing.
	time.Sleep(100 * time.Millisecond)

	log.Println("Shutdown complete")
	return nil
}
