package handler

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Identity.API ok", "data": nil})
}
