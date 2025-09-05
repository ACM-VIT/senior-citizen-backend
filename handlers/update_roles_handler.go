package handlers

import (
	"acm-portal/controllers"
	requestschemas "acm-portal/schemas/request_schemas"

	"github.com/gofiber/fiber/v2"
)

func UpdateRolesHandler(c *fiber.Ctx) error {
	var payload []requestschemas.UpdateRolesRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	results, err := controllers.UpdateRoles(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Roles updated successfully", "data": results})
}
