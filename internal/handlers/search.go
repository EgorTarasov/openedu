package handlers

import (
	"openedu/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) Search(c *fiber.Ctx) error {
	problemID := c.Query("p")
	query := c.Query("q")

	if problemID == "" && query == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var results []models.DBProblem

	if problemID != "" {
		if err := h.db.Where("problem_id LIKE ?", "problem_"+problemID).Find(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to query database")
		}
	} else if query != "" {
		if err := h.db.Raw(`
            SELECT *, 
                   similarity(question, ?) as match_score
            FROM db_problems
            WHERE 
                -- Trigram similarity match (more flexible)
                question % ? OR
                -- ILIKE for substring matching
                question ILIKE ? OR
                -- Fuzzy word matching
                to_tsvector('russian', question) @@ plainto_tsquery('russian', ?)
            ORDER BY match_score DESC
            LIMIT 20
        `, query, query, "%"+query+"%", query).Find(&results).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to perform similarity search")
		}
	}

	return c.JSON(results)
}
