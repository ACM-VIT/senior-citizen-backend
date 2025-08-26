package handlers

import (
	"acm-portal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UpdateRolesHandler(c *fiber.Ctx) error {
	var payload []struct {
		Email    string `json:"email"`
		Role     string `json:"role"`
		Position string `json:"position,omitempty"`
		Tenure   string `json:"tenure"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := controllers.UpdateRoles(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "memberships updated successfully",
	})
}
