package controllers

import (
	"acm-portal/database"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateAdmin(c *fiber.Ctx) error {
	currentUserID, ok := c.Locals("userID").(string)
	if !ok || currentUserID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or missing userID")
	}

	type createAdminRequest struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	}

	var body createAdminRequest
	if err := c.BodyParser(&body); err != nil && string(c.Body()) != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	targetUserID := body.UserID
	if targetUserID == "" && body.Email == "" {
		targetUserID = currentUserID
	}

	if targetUserID == "" && body.Email != "" {
		var user models.User
		if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		targetUserID = user.ID
	}

	var existing models.Admin
	if err := database.DB.Where("user_id = ?", targetUserID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Admin already exists for this user"})
	}

	admin := models.Admin{UserID: targetUserID}
	if err := database.DB.Create(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(admin)
}
