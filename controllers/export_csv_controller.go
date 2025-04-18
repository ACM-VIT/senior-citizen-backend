package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"bytes"
	"encoding/csv"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

func sanitize(input string) string {
	if !utf8.ValidString(input) {
		return string([]rune(input))
	}
	return input
}

func ExportMembersCSV(c *fiber.Ctx) error {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve members for CSV export.")
	}

	if len(users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No members found to export.")
	}

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF}) 
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{"Name", "Email"}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to write CSV header.")
	}

	for _, u := range users {
		if err := writer.Write([]string{
			sanitize(u.Name),
			sanitize(u.Email),
		}); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to write CSV data.")
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate CSV.")
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=members.csv")

	return c.Send(buf.Bytes())
}
