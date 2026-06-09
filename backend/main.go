package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"indo-stats-backend/cache"
	"indo-stats-backend/config"
	"indo-stats-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var startTime time.Time

// main is the entry point of the application.
func main() {
	startTime = time.Now()
	// Load env
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize Cache
	cache.Init()
	log.Println("Cache initialized")

	// Ambil PORT dari env
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Create Fiber app
	app := fiber.New()

	// Middleware Recover
	app.Use(recover.New())

	// Middleware Logger
	app.Use(logger.New())

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Route group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health Check Route
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":            "ok",
			"version":           "1.0.0",
			"uptime":            time.Since(startTime).Round(time.Second).String(),
			"cache_items_count": cache.ItemCount(),
			"endpoints":         []string{"rates", "weather", "commodities", "market", "quakes"},
		})
	})

	// Analytics Routes
	v1.Get("/rates", config.StrictLimiter, handlers.GetRates)
	v1.Get("/weather", config.NormalLimiter, handlers.GetWeather)
	v1.Get("/commodities", config.StrictLimiter, handlers.GetCommodities)
	v1.Get("/market", config.StrictLimiter, handlers.GetMarket)
	v1.Get("/quakes", config.NormalLimiter, handlers.GetQuakes)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server Listen ended: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}