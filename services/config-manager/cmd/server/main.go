package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/historian/config-manager/internal/api"
	"github.com/historian/config-manager/internal/generator"
	"github.com/historian/config-manager/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "historian")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "historian_config")
	configDir := getEnv("CONFIG_DIR", "./config/generated")
	port := getEnv("PORT", "8090")

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to database")

	// Initialize repositories and services
	deviceRepo := repository.NewDeviceRepository(db)
	configGen := generator.NewConfigGenerator(configDir)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "config-manager",
			"version": "2.0.0",
		})
	})

	// Initialize API handlers
	apiHandler := api.NewHandler(deviceRepo, configGen)

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Device management
		v1.GET("/devices", apiHandler.ListDevices)
		v1.GET("/devices/:id", apiHandler.GetDevice)
		v1.DELETE("/devices/:id", apiHandler.DeleteDevice)
		v1.GET("/devices/:id/config", apiHandler.GetDeviceConfig) // New: Config API
		v1.POST("/devices/:id/deploy", apiHandler.DeployDevice)   // New: Deploy Button
		v1.POST("/devices/:id/stop", apiHandler.StopDevice)       // New: Stop Button

		// Modbus-specific
		// Modbus-specific
		modbus := v1.Group("/devices/modbus")
		{
			modbus.GET("", apiHandler.ListModbusDevices)
			modbus.POST("", apiHandler.CreateModbusDevice)
			modbus.GET("/:id", apiHandler.GetModbusDevice)
			modbus.PUT("/:id", apiHandler.UpdateModbusDevice)
		}

		// OPC UA-specific
		opc := v1.Group("/devices/opc")
		{
			opc.GET("", apiHandler.ListOPCDevices)
			opc.POST("", apiHandler.CreateOPCDevice)
			opc.GET("/:id", apiHandler.GetOPCDevice)
		}

		// Config generation
		config := v1.Group("/config")
		{
			config.POST("/generate/:id", apiHandler.GenerateConfig)
			config.GET("/latest/:id", apiHandler.GetLatestConfig)
		}
	}

	// Start server
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("🚀 Config Manager starting on %s", addr)
	log.Printf("📁 Config directory: %s", configDir)
	log.Printf("📊 API docs: http://localhost:%s/health", port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
