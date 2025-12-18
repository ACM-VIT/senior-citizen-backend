package controllers

import (
	"strconv"
	"strings"
	"time"

	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type eventCreateInput struct {
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	TenureID          string           `json:"tenure_id"`
	Venue             string           `json:"venue"`
	EventType         models.EventType `json:"event_type"`
	Date              time.Time        `json:"date"`
	StartTime         *time.Time       `json:"start_time"`
	EndTime           *time.Time       `json:"end_time"`
	NumberOfAttendees uint             `json:"number_of_attendees"`
	NumberOfWinners   uint             `json:"number_of_winners"`
	SponsorIDs        []string         `json:"sponsor_ids"`
}

type eventUpdateInput struct {
	Name              *string           `json:"name"`
	Description       *string           `json:"description"`
	TenureID          *string           `json:"tenure_id"`
	Venue             *string           `json:"venue"`
	EventType         *models.EventType `json:"event_type"`
	Date              *time.Time        `json:"date"`
	StartTime         *time.Time        `json:"start_time"`
	EndTime           *time.Time        `json:"end_time"`
	NumberOfAttendees *uint             `json:"number_of_attendees"`
	NumberOfWinners   *uint             `json:"number_of_winners"`
	SponsorIDs        []string          `json:"sponsor_ids"`
}

func ListEvents(c *fiber.Ctx) error {
	query := database.DB.Model(&models.Event{})
	countQuery := database.DB.Model(&models.Event{})

	if tenureID := strings.TrimSpace(c.Query("tenure_id")); tenureID != "" {
		query = query.Where("tenure_id = ?", tenureID)
		countQuery = countQuery.Where("tenure_id = ?", tenureID)
	}

	if eventType := strings.TrimSpace(c.Query("event_type")); eventType != "" {
		query = query.Where("event_type = ?", eventType)
		countQuery = countQuery.Where("event_type = ?", eventType)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	page, pageSize, offset := helpers.NormalizePagination(page, pageSize)
	query = query.Offset(offset).Limit(pageSize).Order("date desc")

	var events []models.Event
	if err := query.Preload("Sponsors").Preload("Tenure").Find(&events).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch events.")
	}

	var totalCount int64
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count events.")
	}

	_, totalPages := helpers.CalculatePagination(page, pageSize, totalCount)

	return c.JSON(fiber.Map{
		"data":        events,
		"total":       totalCount,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

func GetEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event id is required")
	}

	var event models.Event
	if err := database.DB.Preload("Sponsors").Preload("Tenure").First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Event not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch event.")
	}

	return c.JSON(event)
}

func CreateEvent(c *fiber.Ctx) error {
	var body eventCreateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	if strings.TrimSpace(body.Name) == "" || strings.TrimSpace(body.TenureID) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name and tenure_id are required")
	}

	if !isValidEventType(body.EventType) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid event_type")
	}

	if err := database.DB.First(&models.Tenure{}, "id = ?", body.TenureID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Tenure not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch tenure.")
	}

	var sponsors []models.Sponsor
	if body.SponsorIDs != nil {
		if err := database.DB.Where("id IN ?", body.SponsorIDs).Find(&sponsors).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch sponsors.")
		}
		if len(sponsors) != len(body.SponsorIDs) {
			return fiber.NewError(fiber.StatusBadRequest, "one or more sponsor_ids are invalid")
		}
	}

	event := models.Event{
		Name:              strings.TrimSpace(body.Name),
		Description:       body.Description,
		TenureID:          body.TenureID,
		Venue:             body.Venue,
		EventType:         body.EventType,
		Date:              body.Date,
		StartTime:         body.StartTime,
		EndTime:           body.EndTime,
		NumberOfAttendees: body.NumberOfAttendees,
		NumberOfWinners:   body.NumberOfWinners,
		Sponsors:          sponsors,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.DB.Preload("Sponsors").Preload("Tenure").First(&event, "id = ?", event.ID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch created event.")
	}

	return c.Status(fiber.StatusCreated).JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event id is required")
	}

	var body eventUpdateInput
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body provided.")
	}

	var event models.Event
	if err := database.DB.First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Event not found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch event.")
	}

	updates := map[string]interface{}{}
	if body.Name != nil {
		trimmed := strings.TrimSpace(*body.Name)
		if trimmed == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		updates["name"] = trimmed
	}
	if body.Description != nil {
		updates["description"] = *body.Description
	}
	if body.TenureID != nil {
		if err := database.DB.First(&models.Tenure{}, "id = ?", *body.TenureID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fiber.NewError(fiber.StatusNotFound, "Tenure not found.")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch tenure.")
		}
		updates["tenure_id"] = *body.TenureID
	}
	if body.Venue != nil {
		updates["venue"] = *body.Venue
	}
	if body.EventType != nil {
		if !isValidEventType(*body.EventType) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid event_type")
		}
		updates["event_type"] = *body.EventType
	}
	if body.Date != nil {
		updates["date"] = *body.Date
	}
	if body.StartTime != nil {
		updates["start_time"] = *body.StartTime
	}
	if body.EndTime != nil {
		updates["end_time"] = *body.EndTime
	}
	if body.NumberOfAttendees != nil {
		updates["number_of_attendees"] = *body.NumberOfAttendees
	}
	if body.NumberOfWinners != nil {
		updates["number_of_winners"] = *body.NumberOfWinners
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&event).Updates(updates).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update event.")
		}
	}

	if body.SponsorIDs != nil {
		var sponsors []models.Sponsor
		if err := database.DB.Where("id IN ?", body.SponsorIDs).Find(&sponsors).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch sponsors.")
		}
		if len(sponsors) != len(body.SponsorIDs) {
			return fiber.NewError(fiber.StatusBadRequest, "one or more sponsor_ids are invalid")
		}
		if err := database.DB.Model(&event).Association("Sponsors").Replace(sponsors); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update sponsors.")
		}
	}

	if err := database.DB.Preload("Sponsors").Preload("Tenure").First(&event, "id = ?", event.ID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch updated event.")
	}

	return c.JSON(event)
}

func DeleteEvent(c *fiber.Ctx) error {
	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event id is required")
	}

	if err := database.DB.Delete(&models.Event{}, "id = ?", eventID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete event.")
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Event deleted successfully"})
}

func isValidEventType(eventType models.EventType) bool {
	switch eventType {
	case models.Online, models.Offline, models.Hybrid:
		return true
	default:
		return false
	}
}
