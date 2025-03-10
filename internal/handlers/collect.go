package handlers

import (
	"log/slog"
	"openedu/internal/models"
	"openedu/internal/parser"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) CollectHandler(c *fiber.Ctx) error {
	var payload models.Payload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// TODO: create parsing function
	go func() {
		h.db.Create(&models.DBPayload{Payload: payload})
		problems := parser.ParseContent(payload.Data)
		for _, v := range problems {
			p := models.FromProblem(v)
			h.db.Save(&p)
		}
		slog.Info("saved all problems")
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data received successfully",
	})
}
