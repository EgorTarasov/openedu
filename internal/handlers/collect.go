package handlers

import (
	"openedu/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) CollectHandler(c *fiber.Ctx) error {
	var payload models.Payload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	h.db.Create(&models.DBPayload{Payload: payload})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data received successfully",
	})
}
