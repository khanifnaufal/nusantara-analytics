package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"indo-stats-backend/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
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

	// Parse RATE_LIMIT_MAX
	rateLimitMax := 60
	if maxStr := os.Getenv("RATE_LIMIT_MAX"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil {
			rateLimitMax = val
		}
	}

	// Parse RATE_LIMIT_WINDOW_SEC
	rateLimitWindowSec := 60
	if windowStr := os.Getenv("RATE_LIMIT_WINDOW_SEC"); windowStr != "" {
		if val, err := strconv.Atoi(windowStr); err == nil {
			rateLimitWindowSec = val
		}
	}

	// Create Fiber app
	app := fiber.New()

	// Middleware Recover
	app.Use(recover.New())

	// Middleware Logger
	app.Use(logger.New())

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Middleware Limiter
	app.Use(limiter.New(limiter.Config{
		Max:        rateLimitMax,
		Expiration: time.Duration(rateLimitWindowSec) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "too many requests",
				"retry_after": strconv.Itoa(rateLimitWindowSec) + "s",
			})
		},
	}))

	// Route group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health Check Route
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Start server
	log.Fatal(app.Listen(":" + port))
}