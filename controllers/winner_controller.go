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

type winnerCreateInput struct {
	ParticipationID string `json:"participation_id"`
	Position        uint   `json:"position"`
}

type winnerUpdateInput struct {
	ParticipationID *string `json:"participation_id"`
	Position        *uint   `json:"position"`
}

func ListWinners(c *fiber.Ctx) error {
	query := database.DB.Model(&models.Winner{}).
		Preload("Participation.User.Profile").
		Preload("Participation.Event.Tenure")
	countQuery := database.DB.Model(&models.Winner{})

	if eventID := strings.TrimSpace(c.Query("event_id")); eventID != "" {
		query = query.Where("event_id = ?", eventID)
		countQuery = countQuery.Where("event_id = ?", eventID)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	page, pageSize, offset := helpers.NormalizePagination(page, pageSize)
	query = query.Offset(offset).Limit(pageSize).Order("position asc")

	var winners []models.Winner
	if err := query.Find(&winners).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch winners.")
	}

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count winners.")
	}

	_, totalPages := helpers.CalculatePagination(page, pageSize, totalCount)

	return c.JSON(fiber.Map{
		"data":        winners,
		"total":       totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetWinner(c *fiber.Ctx) error {
	winnerID := c.Params("id")
	if winnerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "winner id is required")
	}

	var winner models.Winner
	if err := database.DB.Preload("Participation.User.Profile").Preload("Participation.Event.Tenure").First(&winner, "id = ?", winnerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Winner not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch winner.")
	}

	return c.JSON(winner)
}

func CreateWinner(c *fiber.Ctx) error {
	var body winnerCreateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	if strings.TrimSpace(body.ParticipationID) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "participation_id is required")
	}
	if body.Position == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "position is required")
	}

	var participation models.Participation
	if err := database.DB.First(&participation, "id = ?", body.ParticipationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Participation not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch participation.")
	}

	var existing models.Winner
	if err := database.DB.Where("event_id = ? AND position = ?", participation.EventID, body.Position).First(&existing).Error; err == nil {
		return fiber.NewError(fiber.StatusConflict, "Winner position already assigned for this event.")
	} else if err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check winner position.")
	}

	winner := models.Winner{
		ParticipationID: body.ParticipationID,
		EventID:         participation.EventID,
		Position:        body.Position,
	}

	if err := database.DB.Create(&winner).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.DB.Preload("Participation.User.Profile").Preload("Participation.Event.Tenure").First(&winner, "id = ?", winner.ID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch created winner.")
	}

	return c.Status(fiber.StatusCreated).JSON(winner)
}

func UpdateWinner(c *fiber.Ctx) error {
	winnerID := c.Params("id")
	if winnerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "winner id is required")
	}

	var body winnerUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	var winner models.Winner
	if err := database.DB.First(&winner, "id = ?", winnerID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Winner not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch winner.")
	}

	updates := map[string]interface{}{}
	if body.ParticipationID != nil {
		var participation models.Participation
		if err := database.DB.First(&participation, "id = ?", *body.ParticipationID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fiber.NewError(fiber.StatusNotFound, "Participation not found.")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch participation.")
		}
		updates["participation_id"] = *body.ParticipationID
		updates["event_id"] = participation.EventID
	}
	if body.Position != nil {
		if *body.Position == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "position must be greater than zero")
		}
		updates["position"] = *body.Position
	}

	if len(updates) > 0 {
		eventID := winner.EventID
		if updatedEventID, ok := updates["event_id"].(string); ok && updatedEventID != "" {
			eventID = updatedEventID
		}
		position := winner.Position
		if updatedPosition, ok := updates["position"].(uint); ok && updatedPosition > 0 {
			position = updatedPosition
		}

		var existing models.Winner
		if err := database.DB.Where("id <> ? AND event_id = ? AND position = ?", winnerID, eventID, position).First(&existing).Error; err == nil {
			return fiber.NewError(fiber.StatusConflict, "Winner position already assigned for this event.")
		} else if err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to check winner position.")
		}

		if err := database.DB.Model(&winner).Updates(updates).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update winner.")
		}
	}

	if err := database.DB.Preload("Participation.User.Profile").Preload("Participation.Event.Tenure").First(&winner, "id = ?", winnerID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated winner.")
	}

	return c.JSON(winner)
}

func DeleteWinner(c *fiber.Ctx) error {
	winnerID := c.Params("id")
	if winnerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "winner id is required")
	}

	if err := database.DB.Delete(&models.Winner{}, "id = ?", winnerID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete winner.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Winner deleted successfully"})
}
