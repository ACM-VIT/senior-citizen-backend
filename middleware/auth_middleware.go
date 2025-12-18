package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"acm-portal/controllers"
	"acm-portal/internal/config"

	"github.com/gofiber/fiber/v2"
)

type HydraIntrospectionResponse struct {
	Active   bool     `json:"active"`
	Sub      string   `json:"sub"`
	Exp      int64    `json:"exp"`
	Scope    string   `json:"scope"`
	ClientID string   `json:"client_id"`
	Email    string   `json:"email,omitempty"`
	Audience []string `json:"aud,omitempty"`
}

func extractToken(c *fiber.Ctx) string {
	bearerToken := strings.TrimSpace(c.Get("Authorization"))
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer "))
	}
	return ""
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := extractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token required"})
	}

	cfg := config.Get()
	if cfg.HydraAdminURL == "" || cfg.KratosAdminURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "missing auth service configuration"})
	}

	hydraAdminURL := strings.TrimSuffix(cfg.HydraAdminURL, "/")
	introspectURL := hydraAdminURL + "/admin/oauth2/introspect"
	resp, err := http.PostForm(introspectURL, url.Values{
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid response from introspection endpoint"})
	}

	var introspection HydraIntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode introspection response", "details": err.Error()})
	}

	if !introspection.Active || introspection.Exp < time.Now().Unix() || introspection.Sub == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Inactive or expired token"})
	}
	kratosAdminURL := strings.TrimSuffix(cfg.KratosAdminURL, "/")
	kratosResp, err := http.Get(fmt.Sprintf("%s/admin/identities/%s", kratosAdminURL, introspection.Sub))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user info from Kratos",
			"details": err.Error(),
		})
	}
	defer kratosResp.Body.Close()

	if kratosResp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch identity: status %d", kratosResp.StatusCode)})
	}

	var kratosIdentity struct {
		ID     string                 `json:"id"`
		Traits map[string]interface{} `json:"traits"`
	}
	if err := json.NewDecoder(kratosResp.Body).Decode(&kratosIdentity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to decode Kratos identity",
			"details": err.Error(),
		})
	}

	email, _ := kratosIdentity.Traits["email"].(string)
	if strings.TrimSpace(email) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing email in identity traits"})
	}

	userID, err := controllers.GetOrCreateUser(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to ensure user exists"})
	}

	claims := fiber.Map{
		"sub":      introspection.Sub,
		"email":    email,
		"scope":    introspection.Scope,
		"clientID": introspection.ClientID,
		"exp":      introspection.Exp,
	}

	c.Locals("claims", claims)
	c.Locals("userID", userID)
	c.Locals("kratosID", introspection.Sub)

	return c.Next()
}
