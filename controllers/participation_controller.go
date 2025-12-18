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

type participationCreateInput struct {
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	TeamName string `json:"team_name"`
}

type participationUpdateInput struct {
	EventID  *string `json:"event_id"`
	UserID   *string `json:"user_id"`
	TeamName *string `json:"team_name"`
}

func ListParticipations(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	isAdmin, err := helpers.IsAdmin(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check admin status.")
	}

	query := database.DB.Model(&models.Participation{}).
		Preload("User.Profile").
		Preload("Event.Tenure")
	countQuery := database.DB.Model(&models.Participation{})

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
		countQuery = countQuery.Where("user_id = ?", userID)
	}

	if eventID := strings.TrimSpace(c.Query("event_id")); eventID != "" {
		query = query.Where("event_id = ?", eventID)
		countQuery = countQuery.Where("event_id = ?", eventID)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	page, pageSize, offset := helpers.NormalizePagination(page, pageSize)
	query = query.Offset(offset).Limit(pageSize).Order("created_at desc")

	var participations []models.Participation
	if err := query.Find(&participations).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch participations.")
	}

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count participations.")
	}

	_, totalPages := helpers.CalculatePagination(page, pageSize, totalCount)

	return c.JSON(fiber.Map{
		"data":        participations,
		"total":       totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetParticipation(c *fiber.Ctx) error {
	participationID := c.Params("id")
	if participationID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "participation id is required")
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	isAdmin, err := helpers.IsAdmin(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check admin status.")
	}

	var participation models.Participation
	if err := database.DB.Preload("User.Profile").Preload("Event.Tenure").First(&participation, "id = ?", participationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Participation not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch participation.")
	}

	if !isAdmin && participation.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Not authorized to view this participation.")
	}

	return c.JSON(participation)
}

func CreateParticipation(c *fiber.Ctx) error {
	var body participationCreateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	if strings.TrimSpace(body.EventID) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event_id is required")
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	isAdmin, err := helpers.IsAdmin(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check admin status.")
	}

	targetUserID := userID
	if strings.TrimSpace(body.UserID) != "" {
		if !isAdmin && body.UserID != userID {
			return fiber.NewError(fiber.StatusForbidden, "Not authorized to create participation for another user.")
		}
		targetUserID = body.UserID
	}

	if err := database.DB.First(&models.User{}, "id = ?", targetUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "User not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch user.")
	}

	if err := database.DB.First(&models.Event{}, "id = ?", body.EventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Event not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch event.")
	}

	var existing models.Participation
	if err := database.DB.Where("user_id = ? AND event_id = ?", targetUserID, body.EventID).First(&existing).Error; err == nil {
		return fiber.NewError(fiber.StatusConflict, "Participation already exists.")
	} else if err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check existing participation.")
	}

	participation := models.Participation{
		UserID:   targetUserID,
		EventID:  body.EventID,
		TeamName: body.TeamName,
	}

	if err := database.DB.Create(&participation).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.DB.Preload("User.Profile").Preload("Event.Tenure").First(&participation, "id = ?", participation.ID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch created participation.")
	}

	return c.Status(fiber.StatusCreated).JSON(participation)
}

func UpdateParticipation(c *fiber.Ctx) error {
	participationID := c.Params("id")
	if participationID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "participation id is required")
	}

	var body participationUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	isAdmin, err := helpers.IsAdmin(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check admin status.")
	}

	var participation models.Participation
	if err := database.DB.First(&participation, "id = ?", participationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Participation not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch participation.")
	}

	if !isAdmin && participation.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Not authorized to update this participation.")
	}

	updates := map[string]interface{}{}
	if body.TeamName != nil {
		updates["team_name"] = *body.TeamName
	}

	if body.EventID != nil || body.UserID != nil {
		if !isAdmin {
			return fiber.NewError(fiber.StatusForbidden, "Not authorized to change event or user.")
		}
		if body.EventID != nil {
			if err := database.DB.First(&models.Event{}, "id = ?", *body.EventID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return fiber.NewError(fiber.StatusNotFound, "Event not found.")
				}
				return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch event.")
			}
			updates["event_id"] = *body.EventID
		}
		if body.UserID != nil {
			if err := database.DB.First(&models.User{}, "id = ?", *body.UserID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return fiber.NewError(fiber.StatusNotFound, "User not found.")
				}
				return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch user.")
			}
			updates["user_id"] = *body.UserID
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&participation).Updates(updates).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update participation.")
		}
	}

	if err := database.DB.Preload("User.Profile").Preload("Event.Tenure").First(&participation, "id = ?", participation.ID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated participation.")
	}

	return c.JSON(participation)
}

func DeleteParticipation(c *fiber.Ctx) error {
	participationID := c.Params("id")
	if participationID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "participation id is required")
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Token invalid or expired.")
	}

	isAdmin, err := helpers.IsAdmin(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check admin status.")
	}

	var participation models.Participation
	if err := database.DB.First(&participation, "id = ?", participationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Participation not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch participation.")
	}

	if !isAdmin && participation.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "Not authorized to delete this participation.")
	}

	if err := database.DB.Delete(&models.Participation{}, "id = ?", participationID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete participation.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Participation deleted successfully"})
}
