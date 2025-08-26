package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func UpdateInformation(c *fiber.Ctx) error {
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

	var user models.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		log.Println(userID)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user record. Please try again later.")
	}

	var updateData struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Gender string `json:"gender"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided. Please check your request format.")
	}

	if err := tx.Model(&user).Updates(updateData).Error; err != nil {
		log.Println("Error updating user:", err)
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user information.")
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Transaction commit failed:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save user info. Please try again.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Information updated successfully", "data": user})
}
