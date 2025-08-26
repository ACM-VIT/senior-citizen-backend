package handlers

import (
	"acm-portal/controllers"
	"acm-portal/helpers"

	"github.com/gofiber/fiber/v2"
)

// SignupOrLoginHandler - fetch existing user or create new one
func SignupOrLoginHandler(c *fiber.Ctx) error {
	type request struct {
		Email  string `json:"email"`
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Gender string `json:"gender"`
		Admin  bool   `json:"admin"`
	}

	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := controllers.FetchOrCreateUser(body.Email, body.Name, body.Phone, body.Gender)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := helpers.GenerateJWT(userID, body.Email, body.Admin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
