package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CallbackHandler(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing code",
		})
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:3000/callback") // must match client registration
	data.Set("client_id", "test-client2")
	data.Set("client_secret", "test-secret2")

	resp, err := http.Post(
		"http://127.0.0.1:4444/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to exchange code for token",
		})
	}
	defer resp.Body.Close()

	var tokens map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to decode token response",
		})
	}

	// For now just return tokens (later: set cookie/session)
	return c.JSON(tokens)
}
