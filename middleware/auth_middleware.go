package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HydraIntrospectionResponse represents the token introspection result from Hydra
type HydraIntrospectionResponse struct {
	Active   bool     `json:"active"`
	Sub      string   `json:"sub"`
	Exp      int64    `json:"exp"`
	Scope    string   `json:"scope"`
	ClientID string   `json:"client_id"`
	Username string   `json:"username,omitempty"`
	Email    string   `json:"email,omitempty"` // if mapped in claims
	Admin    bool     `json:"admin,omitempty"`
	Audience []string `json:"aud,omitempty"`
}

// extractToken gets the JWT from Authorization header or cookie
func extractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return c.Cookies(os.Getenv("CSRF_COOKIE_NAME"))
}

// AuthMiddleware validates the token via Hydra introspection
func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := extractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	// Call Hydra introspection endpoint
	resp, err := http.PostForm("http://127.0.0.1:4444/oauth2/introspect", url.Values{
		"token": {tokenString},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to introspect token",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid response from introspection endpoint",
		})
	}

	var introspection HydraIntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to decode introspection response",
			"details": err.Error(),
		})
	}

	// Check if token is active
	if !introspection.Active {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Inactive or invalid token",
		})
	}

	// Check expiration
	if introspection.Exp < time.Now().Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expired",
		})
	}

	// Attach user info to request context
	c.Locals("userID", introspection.Sub)
	c.Locals("email", introspection.Email)
	c.Locals("scopes", introspection.Scope)
	c.Locals("clientID", introspection.ClientID)
	c.Locals("admin", introspection.Admin)

	return c.Next()
}