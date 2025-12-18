package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func CallbackHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code":              c.Query("code"),
		"state":             c.Query("state"),
		"error":             c.Query("error"),
		"error_description": c.Query("error_description"),
	})
}
