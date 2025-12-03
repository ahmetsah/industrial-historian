package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/handler"
	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/middleware"
	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/repository"
	"github.com/ahmetsah/industrial-historian/go-services/auth/internal/service"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	seedAdmin := flag.Bool("seed-admin", false, "Seed initial admin user")
	adminUser := flag.String("admin-user", "admin", "Admin username")
	adminPass := flag.String("admin-pass", "admin123", "Admin password")
	flag.Parse()

	// Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/historian?sslmode=disable"
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "private.pem"
	}

	// Database Connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Services and Repositories
	userRepo := repository.NewPostgresUserRepository(db)

	if *seedAdmin {
		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(*adminPass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		user := &repository.User{
			Username:     *adminUser,
			PasswordHash: string(hash),
			Role:         "ADMIN",
		}
		if err := userRepo.CreateUser(user); err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
		log.Printf("Admin user '%s' created successfully", *adminUser)
		return
	}

	// NATS Connection
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	tokenService, err := service.NewTokenService(privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to initialize token service: %v", err)
	}

	authHandler := &handler.AuthHandler{
		Repo:         userRepo,
		TokenService: tokenService,
		NatsConn:     nc,
	}

	// Middleware
	rbac := middleware.NewRBACMiddleware(tokenService)

	// Routes
	http.HandleFunc("/api/v1/login", authHandler.Login)
	http.HandleFunc("/api/v1/re-auth", authHandler.ReAuth)
	
	// Protected Routes
	http.Handle("/api/v1/service-accounts", rbac.Authenticate(
		rbac.RequireRole("ADMIN")(http.HandlerFunc(authHandler.CreateServiceAccount)),
	))

	// Start Server
	server := &http.Server{
		Addr: ":" + port,
	}

	go func() {
		log.Printf("Starting auth service on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
