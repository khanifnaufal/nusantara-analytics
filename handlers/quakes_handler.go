package handlers

import (
	"indo-stats-backend/services"

	"github.com/gofiber/fiber/v2"
)

func GetQuakes(c *fiber.Ctx) error {
	resp, fromCache, err := services.FetchQuakes()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": err.Error(),
			"code":  503,
		})
	}

	if fromCache {
		c.Set("X-Cache", "HIT")
	} else {
		c.Set("X-Cache", "MISS")
	}

	return c.JSON(resp)
}
