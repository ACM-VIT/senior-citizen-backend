package middleware

import (

	"acm-portal/database"
	"acm-portal/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		userId := user.ID
		log.Printf("Checking if user %s is an admin", userId)
		var admin models.Admin
		if err := database.DB.First(&admin, "user_id = ?", userId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("User %s is not an admin", userId)
				return &fiber.Error{Code: 401, Message: "Unauthorized admin access"}
			}
			log.Printf("Error fetching admin record for user %s: %v", userId, err)
			return &fiber.Error{Code: 500, Message: "Database error"}
		}
		log.Printf("User %s is an admin", userId)
		c.Locals("admin", admin)
		return c.Next()
	}
}
