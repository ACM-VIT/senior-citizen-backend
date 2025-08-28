package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	/*return c.JSON(fiber.Map{
		"message": "Access granted",
		"userID":  userID,
	})*/

	if userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}
	//log.Println(userID)
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Transaction Failed: Recovered from panic:", r)
			tx.Rollback()
		}
	}()

	var profile models.Profile
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch profile. Please try again later.")
	}

	var updateData models.Profile
	if err := c.BodyParser(&updateData); err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided. Please check your request format.")
	}

	if err := tx.Model(&profile).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile.")
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Transaction commit failed:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save user profile. Please try again.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Information updated successfully", "data": profile})
}
