package controllers

import (
	"acm-portal/database"
	"acm-portal/models"
	requestschemas "acm-portal/schemas/request_schemas"
	responseschemas "acm-portal/schemas/response_schemas"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetMembers(c *fiber.Ctx) error {
    // adminAny := c.Locals("admin")
    // if adminAny == nil {
    //     return &fiber.Error{Code: 401, Message: "Unauthorized access"}
    // }
    // admin, ok := adminAny.(models.Admin)
    // if !ok {
    //     return &fiber.Error{Code: 500, Message: "Admin context invalid"}
    // }
    // log.Printf("[Admin: %s] Accessed member list\n", admin.UserID)

    var body requestschemas.GetMembersRequest
    if err := c.QueryParser(&body); err != nil {
        return &fiber.Error{Code: 400, Message: "Invalid query parameters"}
    }

    tx := database.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            log.Println("Recovered from panic:", r)
            tx.Rollback()
        }
    }()

    var users []models.User
    query := tx.Model(&models.User{})

    if body.Gender != nil {
        query = query.Where("gender = ?", string(*body.Gender))
    }
    if body.JobTitle != nil {
        query = query.Where("job_title ILIKE ?", "%"+*body.JobTitle+"%")
    }
    if body.Company != nil {
        query = query.Where("company ILIKE ?", "%"+*body.Company+"%")
    }
    if body.CreatedAfter != nil {
        if parsed, err := time.Parse(time.RFC3339, *body.CreatedAfter); err == nil {
            query = query.Where("created_at >= ?", parsed)
        }
    }
    if body.UpdatedAfter != nil {
        if parsed, err := time.Parse(time.RFC3339, *body.UpdatedAfter); err == nil {
            query = query.Where("updated_at >= ?", parsed)
        }
    }
    

    page := 1
    pageSize := 10 
    if body.Page != nil && *body.Page > 0 {
        page = *body.Page
    }
    if body.PageSize != nil && *body.PageSize > 0 {
        pageSize = *body.PageSize
    }
    offset := (page - 1) * pageSize

    var total int64
    if err := query.Count(&total).Error; err != nil {
        tx.Rollback()
        log.Println("Count error:", err)
        return &fiber.Error{Code: 500, Message: "Failed to count members"}
    }
    
    query = query.Offset(offset).Limit(pageSize)
    if err := query.Find(&users).Error; err != nil {
        tx.Rollback()
        log.Println("Error fetching members:", err)
        return &fiber.Error{Code: 500, Message: "Failed to retrieve members"}
    }

    tx.Commit()

    var res []responseschemas.GetMembersResponse
    for _, u := range users {
        res = append(res, responseschemas.GetMembersResponse{
            ID:                u.ID,
            Name:              u.Name,
            Email:             u.Email,
            Phone:             u.Phone,
            Gender:            u.Gender,
            ProfilePictureURL: u.ProfilePictureURL,
        })
    }
    totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
    return c.JSON(responseschemas.PaginatedMembersResponse{
        Data:       res,
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    })
}
