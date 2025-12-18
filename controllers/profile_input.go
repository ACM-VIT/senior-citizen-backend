package controllers

import (
	"strings"

	"acm-portal/models"

	"github.com/gofiber/fiber/v2"
)

type profileUpdateInput struct {
	Name              *string        `json:"name"`
	ProfilePictureURL *string        `json:"profile_picture_url"`
	Phone             *string        `json:"phone"`
	Gender            *models.Gender `json:"gender"`
	Bio               *string        `json:"bio"`
	TheSentence       *string        `json:"the_sentence"`
	LinkedIn          *string        `json:"linkedin"`
	GitHub            *string        `json:"github"`
	Twitter           *string        `json:"twitter"`
	Instagram         *string        `json:"instagram"`
	JobTitle          *string        `json:"job_title"`
	Company           *string        `json:"company"`
	Website           *string        `json:"website"`
}

func buildProfileUpdates(input profileUpdateInput) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		if trimmed == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		updates["name"] = trimmed
	}
	if input.ProfilePictureURL != nil {
		updates["profile_picture_url"] = strings.TrimSpace(*input.ProfilePictureURL)
	}
	if input.Phone != nil {
		trimmed := strings.TrimSpace(*input.Phone)
		if trimmed == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "phone cannot be empty")
		}
		updates["phone"] = trimmed
	}
	if input.Gender != nil {
		if *input.Gender != models.Male && *input.Gender != models.Female && *input.Gender != models.Other {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid gender")
		}
		updates["gender"] = *input.Gender
	}
	if input.Bio != nil {
		updates["bio"] = strings.TrimSpace(*input.Bio)
	}
	if input.TheSentence != nil {
		updates["the_sentence"] = strings.TrimSpace(*input.TheSentence)
	}
	if input.LinkedIn != nil {
		updates["linked_in"] = strings.TrimSpace(*input.LinkedIn)
	}
	if input.GitHub != nil {
		updates["git_hub"] = strings.TrimSpace(*input.GitHub)
	}
	if input.Twitter != nil {
		updates["twitter"] = strings.TrimSpace(*input.Twitter)
	}
	if input.Instagram != nil {
		updates["instagram"] = strings.TrimSpace(*input.Instagram)
	}
	if input.JobTitle != nil {
		updates["job_title"] = strings.TrimSpace(*input.JobTitle)
	}
	if input.Company != nil {
		updates["company"] = strings.TrimSpace(*input.Company)
	}
	if input.Website != nil {
		updates["website"] = strings.TrimSpace(*input.Website)
	}

	return updates, nil
}
