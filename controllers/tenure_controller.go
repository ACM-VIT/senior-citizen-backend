package controllers

import (
	"strconv"
	"time"

	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type tenureCreateInput struct {
	IsCurrentTenure bool      `json:"is_current_tenure"`
	CurrentYear     int       `json:"current_year"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

type tenureUpdateInput struct {
	IsCurrentTenure *bool      `json:"is_current_tenure"`
	CurrentYear     *int       `json:"current_year"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
}

func ListTenures(c *fiber.Ctx) error {
	query := database.DB.Model(&models.Tenure{})
	countQuery := database.DB.Model(&models.Tenure{})

	if c.Query("current") == "true" {
		query = query.Where("is_current_tenure = ?", true)
		countQuery = countQuery.Where("is_current_tenure = ?", true)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	page, pageSize, offset := helpers.NormalizePagination(page, pageSize)
	query = query.Offset(offset).Limit(pageSize).Order("current_year desc")

	var tenures []models.Tenure
	if err := query.Find(&tenures).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch tenures.")
	}

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count tenures.")
	}

	_, totalPages := helpers.CalculatePagination(page, pageSize, totalCount)

	return c.JSON(fiber.Map{
		"data":        tenures,
		"total":       totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetTenure(c *fiber.Ctx) error {
	tenureID := c.Params("id")
	if tenureID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tenure id is required")
	}

	var tenure models.Tenure
	if err := database.DB.First(&tenure, "id = ?", tenureID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Tenure not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch tenure.")
	}

	return c.JSON(tenure)
}

func CreateTenure(c *fiber.Ctx) error {
	var body tenureCreateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	if body.CurrentYear <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "current_year is required")
	}

	tenure := models.Tenure{
		IsCurrentTenure: body.IsCurrentTenure,
		CurrentYear:     body.CurrentYear,
		StartDate:       body.StartDate,
		EndDate:         body.EndDate,
	}

	tx := database.DB.Begin()
	if body.IsCurrentTenure {
		if err := tx.Model(&models.Tenure{}).Where("is_current_tenure = ?", true).Update("is_current_tenure", false).Error; err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update current tenure.")
		}
	}

	if err := tx.Create(&tenure).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save tenure.")
	}

	return c.Status(fiber.StatusCreated).JSON(tenure)
}

func UpdateTenure(c *fiber.Ctx) error {
	tenureID := c.Params("id")
	if tenureID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tenure id is required")
	}

	var body tenureUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	var tenure models.Tenure
	if err := database.DB.First(&tenure, "id = ?", tenureID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Tenure not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch tenure.")
	}

	updates := map[string]interface{}{}
	if body.IsCurrentTenure != nil {
		updates["is_current_tenure"] = *body.IsCurrentTenure
	}
	if body.CurrentYear != nil {
		if *body.CurrentYear <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "current_year must be positive")
		}
		updates["current_year"] = *body.CurrentYear
	}
	if body.StartDate != nil {
		updates["start_date"] = *body.StartDate
	}
	if body.EndDate != nil {
		updates["end_date"] = *body.EndDate
	}

	if len(updates) > 0 {
		tx := database.DB.Begin()
		if body.IsCurrentTenure != nil && *body.IsCurrentTenure {
			if err := tx.Model(&models.Tenure{}).Where("id <> ?", tenureID).Update("is_current_tenure", false).Error; err != nil {
				tx.Rollback()
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to update current tenure.")
			}
		}

		if err := tx.Model(&tenure).Updates(updates).Error; err != nil {
			tx.Rollback()
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update tenure.")
		}

		if err := tx.Commit().Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to save tenure.")
		}
	}

	if err := database.DB.First(&tenure, "id = ?", tenureID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated tenure.")
	}

	return c.JSON(tenure)
}

func DeleteTenure(c *fiber.Ctx) error {
	tenureID := c.Params("id")
	if tenureID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "tenure id is required")
	}

	if err := database.DB.Delete(&models.Tenure{}, "id = ?", tenureID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete tenure.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Tenure deleted successfully"})
}
