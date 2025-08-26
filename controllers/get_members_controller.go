package controllers

import (
	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/models"
	requestschemas "acm-portal/schemas/request_schemas"
	responseschemas "acm-portal/schemas/response_schemas"

	"github.com/gofiber/fiber/v2"
)

func GetMembers(c *fiber.Ctx) error {
	var body requestschemas.GetMembersRequest
	if err := c.QueryParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters provided. Please check your request format.")
	}

	query := database.DB.Model(&models.Membership{})

	if body.Role != "" {
		query = query.Where("role = ?", body.Role)
	}

	if body.Position != "" {
		query = query.Where("position = ?", body.Position)
	}

	if body.Page == 0 || body.Page < 1 {
		body.Page = 1
	}
	if body.PageSize == 0 || body.PageSize < 1 {
		body.PageSize = 10
	}

	offset, _ := helpers.CalculatePagination(body.Page, body.PageSize, 0)
	query = query.Offset(offset).Limit(body.PageSize)

	var memberships []models.Membership
	if err := query.
		Preload("User").
		Preload("Tenure").
		Find(&memberships).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to fetch membership records.")
	}

	if len(memberships) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No members found matching the criteria.")
	}

	var totalCount int64
	if err := database.DB.Model(&models.Membership{}).Count(&totalCount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to count total Membership records. Please try again later.")
	}

	_, totalPages := helpers.CalculatePagination(body.Page, body.PageSize, totalCount)

	var responseData []responseschemas.GetMembersResponse
	for _, membership := range memberships {
		responseData = append(responseData, responseschemas.GetMembersResponse{
			ID:       membership.ID,
			UserID:   membership.UserID,
			User:     membership.User,
			TenureID: membership.TenureID,
			Tenure:   membership.Tenure,
			Role:     membership.Role,
			Position: membership.Position,
		})
	}

	response := responseschemas.PaginatedMembersResponse{
		Data:       responseData,
		Total:      totalCount,
		Page:       body.Page,
		PageSize:   body.PageSize,
		TotalPages: totalPages,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Members fetched successfully",
		"data":    response,
	})
}
