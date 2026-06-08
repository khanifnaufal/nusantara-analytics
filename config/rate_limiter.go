package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// StrictLimiter is rate limiting middleware with a preset of 20 req/minute per IP per route.
// Used for endpoints fetching external APIs: /rates, /commodities, /market.
var StrictLimiter = limiter.New(limiter.Config{
	Max:        20,
	Expiration: 1 * time.Minute,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.IP() + "_" + c.Path()
	},
	LimitReached: func(c *fiber.Ctx) error {
		c.Set("X-Ratelimit-Limit", "20")
		c.Set("X-Ratelimit-Remaining", "0")
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "too many requests",
			"retry_after": "60s",
		})
	},
})

// NormalLimiter is rate limiting middleware with a preset of 30 req/minute per IP per route.
// Used for endpoints fetching external APIs: /weather, /quakes.
var NormalLimiter = limiter.New(limiter.Config{
	Max:        30,
	Expiration: 1 * time.Minute,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.IP() + "_" + c.Path()
	},
	LimitReached: func(c *fiber.Ctx) error {
		c.Set("X-Ratelimit-Limit", "30")
		c.Set("X-Ratelimit-Remaining", "0")
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "too many requests",
			"retry_after": "60s",
		})
	},
})
