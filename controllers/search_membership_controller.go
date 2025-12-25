package controllers

import (
	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	db := database.DB
	var memberships []models.Membership

	if body.Query != "" {
		qGrams := helpers.Trigrams(body.Query)

		var membershipIDs []string

		rawSQL := `
		WITH ranked_profiles AS (
	SELECT
		p.id,
		p.user_id,
		cardinality(
		ARRAY(
			SELECT unnest(COALESCE(p.name_grams, '{}'))
			INTERSECT
			SELECT unnest(?::TEXT[])
		)
		) +
		cardinality(
		ARRAY(
			SELECT unnest(COALESCE(p.company_grams, '{}'))
			INTERSECT
			SELECT unnest(?::TEXT[])
		)
		) AS score
	FROM profiles p
	WHERE
		COALESCE(p.name_grams, '{}') && ?::TEXT[]
		OR COALESCE(p.company_grams, '{}') && ?::TEXT[]
	)
	SELECT m.id
	FROM ranked_profiles rp
	JOIN memberships m ON m.user_id = rp.user_id
	ORDER BY rp.score DESC
	LIMIT 50;
		`

		if err := db.Raw(
			rawSQL,
			pq.Array(qGrams),
			pq.Array(qGrams),
			pq.Array(qGrams),
			pq.Array(qGrams),
		).Scan(&membershipIDs).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if len(membershipIDs) == 0 {
			return c.JSON([]models.Membership{})
		}

		if err := db.
			Where("id IN ?", membershipIDs).
			Preload("User.Profile").
			Preload("Tenure").
			Find(&memberships).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(memberships)
	}

	query := db.
		Model(&models.Membership{}).
		Preload("User.Profile").
		Preload("Tenure")

	if body.TenureID != "" {
		query = query.Where("tenure_id = ?", body.TenureID)
	}
	if body.Role != "" {
		query = query.Where("role = ?", body.Role)
	}
	if body.Position != "" {
		query = query.Where("position = ?", body.Position)
	}

	if err := query.Find(&memberships).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(memberships)
}
