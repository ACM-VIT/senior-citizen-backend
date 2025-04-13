package controllers

import (
	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"
	requestschemas "acm-portal/schemas/request_schemas"
	responseschemas "acm-portal/schemas/response_schemas"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetMembers(c *fiber.Ctx) error {
	var body requestschemas.GetMembersRequest
	if err := c.QueryParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters provided. Please check your request format.")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Transaction Failed: Recovered from panic:", r)
			tx.Rollback()
		}
	}()

	query := tx.Model(&models.Membership{})

	if body.Role != nil {
		query = query.Where("role = ?", *body.Role)
	}

	if body.Position != nil {
		query = query.Where("position = ?", *body.Position)
	}

	if body.Page == nil || *body.Page < 1 {
		*body.Page = 1
	}
	if body.PageSize == nil || *body.PageSize < 1 {
		*body.PageSize = 10
	}

	offset, totalPages := helpers.CalculatePagination(*body.Page, *body.PageSize, 0)
	query = query.Offset(offset).Limit(*body.PageSize)


    var memberships []models.Membership
	if err := tx.Debug().Model(&models.Membership{}).
    Preload("User").Preload("Tenure").
    Find(&memberships).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch membership records.")
    }

	if len(memberships) == 0 {
		tx.Rollback()
		return fiber.NewError(fiber.StatusNotFound, "No members found matching the criteria.")
	}

	var totalCount int64
	if err := tx.Model(&models.Membership{}).Count(&totalCount).Error; err != nil {
		tx.Rollback()
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count total Membership records. Please try again later.")
	}	

	_, totalPages = helpers.CalculatePagination(*body.Page, *body.PageSize, totalCount)

	var responseData []responseschemas.GetMembersResponse
	for _, membership := range memberships {
		responseData = append(responseData, responseschemas.GetMembersResponse{
			ID:        membership.ID,
			UserID:    membership.UserID,
			User:      membership.User,
			TenureID:  membership.TenureID,
			Tenure:    membership.Tenure,
			Role:      membership.Role,
			Position:  membership.Position,
		})
	}

	response := responseschemas.PaginatedMembersResponse{
		Data:       responseData,
		Total:      totalCount,
		Page:       *body.Page,
		PageSize:   *body.PageSize,
		TotalPages: totalPages,
	}

	tx.Commit()
	return c.JSON(fiber.Map{"status":  "success","message": "Members fetched successfully","data": response,})
}