package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	// Ambil PORT dari env
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Create Fiber app
	app := fiber.New()

	// Middleware Logger
	app.Use(logger.New())

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
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