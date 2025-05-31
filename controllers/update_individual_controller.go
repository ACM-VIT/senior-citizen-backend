package controllers

import (
	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"
	requestschemas "acm-portal/schemas/request_schemas"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

var validate = validator.New()

func UpsertIndividual(c *fiber.Ctx) error {
	var input requestschemas.UpdateIndividualInput

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	if err := validate.Struct(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var userID string
	var err error
	if input.UserID != "" {
		userID = input.UserID
	} else {
		userID, err = fetchOrCreateUser(input.Email, input.Name, input.Phone, input.Gender)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Could not resolve user")
		}
	}

	var tenureID string
	if input.TenureID != "" {
		tenureID = input.TenureID
	} else if input.CurrentYear != 0 {
		tenureID, err = fetchTenureByYear(input.CurrentYear)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "Either TenureID or CurrentYear is required")
	}

	err = helpers.ValidateRole(input.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid role provided")
	}


	err = helpers.ValidatePosition(input.Position)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid position provided")
	}

	membership := models.Membership{
		UserID:   userID,
		TenureID: tenureID,
		Role:     models.Role(input.Role),
		Position: input.Position,
	}

	err = database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "tenure_id"}},
		UpdateAll: true,
	}).Create(&membership).Error

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not update individual")
	}

	return c.JSON(fiber.Map{"status":  "success","message": "Individual updated successfully","data":input})
}