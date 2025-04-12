package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	"bytes"
	"encoding/csv"
	"log"
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
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			tx.Rollback()
		}
	}()

	var users []models.User
	if err := tx.Model(&models.User{}).Find(&users).Error; err != nil {
		tx.Rollback()
		log.Println("Error fetching members:", err)
		return &fiber.Error{Code: 500, Message: "Failed to retrieve members"}
	}
	tx.Commit()

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF}) 
	writer := csv.NewWriter(&buf)

	writer.Write([]string{"Name", "Email"})

	for _, u := range users {
		writer.Write([]string{
			sanitize(u.Name),
			sanitize(u.Email),
		})
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Println("Error writing CSV:", err)
		return &fiber.Error{Code: 500, Message: "Failed to generate CSV"}
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=members.csv")

	return c.Send(buf.Bytes())
}