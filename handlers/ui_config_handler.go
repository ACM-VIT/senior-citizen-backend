package handlers

import (
	"encoding/json"
	"strings"

	"acm-portal/internal/config"

	"github.com/gofiber/fiber/v2"
)

func UIConfigHandler(c *fiber.Ctx) error {
	cfg := config.Get()

	payload := map[string]string{
		"kratosPublicUrl": strings.TrimSuffix(cfg.KratosBrowserURL, "/"),
		"publicUrl":       strings.TrimSuffix(cfg.PublicURL, "/"),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("window.AUTH_UI_CONFIG={};")
	}

	c.Set(fiber.HeaderContentType, "application/javascript")
	return c.SendString("window.AUTH_UI_CONFIG=" + string(data) + ";")
}
