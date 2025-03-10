package handlers

import "github.com/gofiber/fiber/v2"

func (h *Handlers) Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}
