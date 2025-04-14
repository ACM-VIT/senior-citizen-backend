package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SearchMembership(c *fiber.Ctx) error {
	query := c.Query("query")
	userID := c.Query("user_id")
	tenureID := c.Query("tenure_id")
	role := c.Query("role")
	position := c.Query("position")

	if query == "" && userID == "" && tenureID == "" && role == "" && position == "" {
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

	// Fuzzy + Exact search
	if query != "" {
		q := "%" + strings.ToLower(query) + "%"
		dbQuery = dbQuery.Where(`
			LOWER(users.name) LIKE ?`,
			q,
		)
	}

	// Optional Filters
	if userID != "" {
		dbQuery = dbQuery.Where("memberships.user_id = ?", userID)
	}
	if tenureID != "" {
		dbQuery = dbQuery.Where("memberships.tenure_id = ?", tenureID)
	}
	if role != "" {
		dbQuery = dbQuery.Where("memberships.role = ?", role)
	}
	if position != "" {
		pq := "%" + strings.ToLower(position) + "%"
		dbQuery = dbQuery.Where("LOWER(memberships.position) LIKE ?", pq)
	}

	if err := dbQuery.Find(&memberships).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch memberships",
		})
	}

	return c.JSON(memberships)
}
