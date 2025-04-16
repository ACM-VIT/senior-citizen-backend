package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func UpdateInformation(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "User ID is required.")
	}

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
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user record. Please try again later.")
	}
	if err := c.BodyParser(&user); err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided. Please check your request format.")
	}

	if err := tx.Model(&user).Updates(user).Error; err != nil {
		log.Println("Error updating user:", err)
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user information.")
	}
	
	tx.Commit()
	
	return c.JSON(fiber.Map{"status":  "success","message": "Information updated successfully","data":user})
}