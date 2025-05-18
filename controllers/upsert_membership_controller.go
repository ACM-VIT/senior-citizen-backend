package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	requestschemas "acm-portal/schemas/request_schemas"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

var validate = validator.New()

func UpsertMembership(c *fiber.Ctx) error {
	var input requestschemas.UpsertMembershipInput

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	if err := validate.Struct(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	membership := models.Membership{
		UserID:   input.UserID,
		TenureID: input.TenureID,
		Role:     models.Role(input.Role),
		Position: input.Position,
	}

	err := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "tenure_id"}},
		UpdateAll: true,
	}).Create(&membership).Error

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not upsert membership")
	}

	return c.JSON(fiber.Map{"status":  "success","message": "Membership upserted successfully","data":input})
}