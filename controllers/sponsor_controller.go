package controllers

import (
	"strconv"
	"strings"

	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type sponsorCreateInput struct {
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
	Website string `json:"website"`
	Contact string `json:"contact"`
}

type sponsorUpdateInput struct {
	Name    *string `json:"name"`
	LogoURL *string `json:"logo_url"`
	Website *string `json:"website"`
	Contact *string `json:"contact"`
}

func ListSponsors(c *fiber.Ctx) error {
	query := database.DB.Model(&models.Sponsor{})
	countQuery := database.DB.Model(&models.Sponsor{})

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	page, pageSize, offset := helpers.NormalizePagination(page, pageSize)
	query = query.Offset(offset).Limit(pageSize).Order("name asc")

	var sponsors []models.Sponsor
	if err := query.Find(&sponsors).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch sponsors.")
	}

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count sponsors.")
	}

	_, totalPages := helpers.CalculatePagination(page, pageSize, totalCount)

	return c.JSON(fiber.Map{
		"data":        sponsors,
		"total":       totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetSponsor(c *fiber.Ctx) error {
	sponsorID := c.Params("id")
	if sponsorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "sponsor id is required")
	}

	var sponsor models.Sponsor
	if err := database.DB.First(&sponsor, "id = ?", sponsorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Sponsor not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch sponsor.")
	}

	return c.JSON(sponsor)
}

func CreateSponsor(c *fiber.Ctx) error {
	var body sponsorCreateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	if strings.TrimSpace(body.Name) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}

	sponsor := models.Sponsor{
		Name:    strings.TrimSpace(body.Name),
		LogoURL: body.LogoURL,
		Website: body.Website,
		Contact: body.Contact,
	}

	if err := database.DB.Create(&sponsor).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(sponsor)
}

func UpdateSponsor(c *fiber.Ctx) error {
	sponsorID := c.Params("id")
	if sponsorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "sponsor id is required")
	}

	var body sponsorUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	var sponsor models.Sponsor
	if err := database.DB.First(&sponsor, "id = ?", sponsorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Sponsor not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch sponsor.")
	}

	updates := map[string]interface{}{}
	if body.Name != nil {
		trimmed := strings.TrimSpace(*body.Name)
		if trimmed == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		updates["name"] = trimmed
	}
	if body.LogoURL != nil {
		updates["logo_url"] = *body.LogoURL
	}
	if body.Website != nil {
		updates["website"] = *body.Website
	}
	if body.Contact != nil {
		updates["contact"] = *body.Contact
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&sponsor).Updates(updates).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update sponsor.")
		}
	}

	if err := database.DB.First(&sponsor, "id = ?", sponsorID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated sponsor.")
	}

	return c.JSON(sponsor)
}

func DeleteSponsor(c *fiber.Ctx) error {
	sponsorID := c.Params("id")
	if sponsorID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "sponsor id is required")
	}

	if err := database.DB.Delete(&models.Sponsor{}, "id = ?", sponsorID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete sponsor.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Sponsor deleted successfully"})
}
