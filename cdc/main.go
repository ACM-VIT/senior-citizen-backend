package main

import (
	"log"

	"github.com/lib/pq"

	"acm-portal/database"
	"acm-portal/helpers"
	"acm-portal/internal/config"

	"github.com/gofiber/fiber/v2"
)

type CDCEvent struct {
	Payload []struct {
		After *struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Company string `json:"company"`
		} `json:"after"`
	} `json:"payload"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	database.ConnectToDB(cfg)
	app := fiber.New()
	updateChan := make(chan CDCEvent, 100)
	go func() {
		for evt := range updateChan {
			tx := database.DB.Begin()
			for _, e := range evt.Payload {
				if e.After == nil {
					continue
				}
				update := map[string]interface{}{
					"name_grams":    pq.Array(helpers.Trigrams(e.After.Name)),
					"company_grams": pq.Array(helpers.Trigrams(e.After.Company)),
				}
				if err := tx.Table("profiles").
					Where("id = ?", e.After.ID).
					Updates(update).Error; err != nil {
					log.Println("CDC update failed for ID", e.After.ID, ":", err)
				}
			}
			tx.Commit()
		}
	}()
	app.Post("/cdc", func(c *fiber.Ctx) error {
		var evt CDCEvent
		if err := c.BodyParser(&evt); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		select {
		case updateChan <- evt:
		default:
			log.Println("CDC event dropped: queue full")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	log.Println("CDC service listening on :8081")
	log.Fatal(app.Listen(":8081"))
}
