package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MembershipSearchRequest struct {
	Query    string `json:"query"`
	TenureID string `json:"tenure_id"`
	Role     string `json:"role"`
	Position string `json:"position"`
}

func SearchMembership(c *fiber.Ctx) error {
	var body MembershipSearchRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.Query == "" && body.TenureID == "" && body.Role == "" && body.Position == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one query/filter is required",
		})
	}

	var memberships []models.Membership
	db := database.DB

	dbQuery := db.Model(&models.Membership{}).
		Preload("User.Profile").
		Preload("Tenure").
		Joins("JOIN users u ON u.id = memberships.user_id").
		Joins("JOIN profiles p ON p.user_id = u.id")

	if body.Query != "" {
		q := strings.ToLower(body.Query)
		dbQuery = dbQuery.Where(`
			LOWER(p.name) ILIKE ? OR
			LOWER(p.phone) ILIKE ? OR
			LOWER(p.gender) ILIKE ? OR
			LOWER(p.job_title) ILIKE ? OR
			LOWER(p.company) ILIKE ? OR
			LOWER(memberships.role) ILIKE ? OR
			LOWER(memberships.position) ILIKE ?`,
			"%"+q+"%", "%"+q+"%", "%"+q+"%",
			"%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%",
		)
	}

	// Filters
	if body.TenureID != "" {
		dbQuery = dbQuery.Where("memberships.tenure_id = ?", body.TenureID)
	}
	if body.Role != "" {
		dbQuery = dbQuery.Where("memberships.role = ?", body.Role)
	}
	if body.Position != "" {
		dbQuery = dbQuery.Where("memberships.position ILIKE ?", "%"+body.Position+"%")
	}

	if err := dbQuery.Find(&memberships).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(memberships)
}
