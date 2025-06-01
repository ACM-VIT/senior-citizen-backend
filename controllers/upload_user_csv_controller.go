package controllers

//this is for uploading the csv, here the user just uploads the csv, and all the values in the csv are stored in the db.

import (
	"acm-portal/database"
	"acm-portal/models"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserCSVController struct{}

func NewUserCSVController() *UserCSVController {
	return &UserCSVController{}
}

type UserCSVRecord struct {
	Name              string
	Email             string
	ProfilePictureURL string
	Phone             string
	Gender            string
	Bio               string
	TheSentence       string
	LinkedIn          string
	GitHub            string
	Twitter           string
	Instagram         string
	JobTitle          string
	Company           string
	Website           string
}

func (uc *UserCSVController) UploadUsersCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	records, err := parseUserCSV(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse CSV: " + err.Error(),
		})
	}

	successCount, errorCount, errors := handleUserCSVUpload(database.DB, records)

	response := fiber.Map{
		"success_count": successCount,
		"error_count":   errorCount,
	}

	if errorCount > 0 {
		errorMessages := make([]string, len(errors))
		for i, e := range errors {
			errorMessages[i] = e.Error()
		}
		response["errors"] = errorMessages
	}

	return c.JSON(response)
}

func parseUserCSV(file *multipart.FileHeader) ([]UserCSVRecord, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records := make([]UserCSVRecord, 0)

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	if err := validateUserCSVHeaders(headers); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) == 0 || (len(record) == 1 && record[0] == "") {
			continue
		}

		userRecord := UserCSVRecord{
			Name:              record[0],
			Email:             record[1],
			ProfilePictureURL: record[2],
			Phone:             record[3],
			Gender:            record[4],
			Bio:               record[5],
			TheSentence:       record[6],
			LinkedIn:          record[7],
			GitHub:            record[8],
			Twitter:           record[9],
			Instagram:         record[10],
			JobTitle:          record[11],
			Company:           record[12],
			Website:           record[13],
		}

		records = append(records, userRecord)
	}

	return records, nil
}

func validateUserCSVHeaders(headers []string) error {
	expectedHeaders := []string{
		"Name", "Email", "ProfilePictureURL", "Phone", "Gender",
		"Bio", "TheSentence", "LinkedIn", "GitHub", "Twitter",
		"Instagram", "JobTitle", "Company", "Website",
	}

	if len(headers) != len(expectedHeaders) {
		return errors.New("invalid number of columns in CSV")
	}

	for i, h := range expectedHeaders {
		if headers[i] != h {
			return errors.New("invalid column header: " + headers[i] + ", expected: " + h)
		}
	}

	return nil
}

func handleUserCSVUpload(db *gorm.DB, records []UserCSVRecord) (int, int, []error) {
	var successCount int
	var errorCount int
	var errorsList []error

	for _, record := range records {
		if record.Name == "" || record.Email == "" || record.Phone == "" || record.Gender == "" {
			errorCount++
			errorsList = append(errorsList, errors.New("missing required fields for record with email: "+record.Email))
			continue
		}

		gender := models.Gender(record.Gender)
		if gender != models.Male && gender != models.Female && gender != models.Other {
			errorCount++
			errorsList = append(errorsList, errors.New("invalid gender for record with email: "+record.Email))
			continue
		}

		user := models.User{
			Name:              record.Name,
			Email:             record.Email,
			ProfilePictureURL: record.ProfilePictureURL,
			Phone:             record.Phone,
			Gender:            gender,
			Bio:               record.Bio,
			TheSentence:       record.TheSentence,
			LinkedIn:          record.LinkedIn,
			GitHub:            record.GitHub,
			Twitter:           record.Twitter,
			Instagram:         record.Instagram,
			JobTitle:          record.JobTitle,
			Company:           record.Company,
			Website:           record.Website,
			Base: models.Base{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		if err := db.Create(&user).Error; err != nil {
			errorCount++
			errorsList = append(errorsList, err)
		} else {
			successCount++
		}
	}

	return successCount, errorCount, errorsList
}
