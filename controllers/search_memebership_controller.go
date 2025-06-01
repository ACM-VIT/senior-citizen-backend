package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MembershipSearchRequest struct {
	Query    string `json:"query"`
	UserID   string `json:"user_id"`
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

	if body.Query == "" && body.UserID == "" && body.TenureID == "" && body.Role == "" && body.Position == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one query/filter is required",
		})
	}

	var memberships []models.Membership
	db := database.DB

	dbQuery := db.
		Preload("User").
		Preload("Tenure").
		Joins("JOIN users ON users.id = memberships.user_id")

	// Fuzzy part
	if body.Query != "" && len(body.Query) >= 3 {
		dbQuery = dbQuery.Where("users.name ILIKE ?", "%"+body.Query+"%")
	}

	// Optional Filters
	if body.UserID != "" {
		dbQuery = dbQuery.Where("memberships.user_id = ?", body.UserID)
	}
	if body.TenureID != "" {
		dbQuery = dbQuery.Where("memberships.tenure_id = ?", body.TenureID)
	}
	if body.Role != "" {
		dbQuery = dbQuery.Where("memberships.role = ?", body.Role)
	}
	if body.Position != "" {
		pq := "%" + strings.ToLower(body.Position) + "%"
		dbQuery = dbQuery.Where("LOWER(memberships.position) LIKE ?", pq)
	}

	if err := dbQuery.Find(&memberships).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch memberships",
		})
	}

	return c.JSON(memberships)
}
