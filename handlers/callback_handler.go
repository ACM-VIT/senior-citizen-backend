package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var clients = map[string]struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}{
	"reverse-coding": {
		ClientID:     "reverse-coding",
		ClientSecret: "rc-secret",
		RedirectURI:  "http://localhost:3000/callback",
	},
	"cryptic-hunt": {
		ClientID:     "cryptic-hunt",
		ClientSecret: "ch-secret",
		RedirectURI:  "http://localhost:3000/callback",
	},
	"c2c": {
		ClientID:     "c2c",
		ClientSecret: "c2c-secret",
		RedirectURI:  "http://localhost:3000/callback",
	},
}

func CallbackHandler(c *fiber.Ctx) error {
	app := c.Query("app")
	if app == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing app identifier"})
	}

	client, exists := clients[app]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unknown app"})
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing code"})
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", client.RedirectURI)
	data.Set("client_id", client.ClientID)
	data.Set("client_secret", client.ClientSecret)

	resp, err := http.Post("http://127.0.0.1:4444/oauth2/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to exchange code for token"})
	}
	defer resp.Body.Close()

	var tokens map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode token response"})
	}

	return c.JSON(tokens)
}
