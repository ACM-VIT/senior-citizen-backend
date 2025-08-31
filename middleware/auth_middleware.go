package middleware

import (
	"acm-portal/controllers"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HydraIntrospectionResponse struct {
	Active   bool     `json:"active"`
	Sub      string   `json:"sub"`
	Exp      int64    `json:"exp"`
	Scope    string   `json:"scope"`
	ClientID string   `json:"client_id"`
	Audience []string `json:"aud,omitempty"`
}

func extractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return c.Cookies(os.Getenv("CSRF_COOKIE_NAME"))
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := extractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	hydraAdminURL := "http://senior-citizen-backend-hydra-1:4445/admin"
	resp, err := http.PostForm(hydraAdminURL+"/oauth2/introspect", url.Values{
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

	if !introspection.Active || introspection.Exp < time.Now().Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Inactive or expired token",
		})
	}

	// Fetch full user info from Kratos using the subject (sub)
	kratosAdminURL := "http://senior-citizen-backend-kratos-1:4434"
	kratosResp, err := http.Get(fmt.Sprintf("%s/admin/identities/%s", kratosAdminURL, introspection.Sub))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch user info from Kratos",
			"details": err.Error(),
		})
	}
	defer kratosResp.Body.Close()

	if kratosResp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch identity: status %d", kratosResp.StatusCode),
		})
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
	name, _ := kratosIdentity.Traits["name"].(string)

	userID, err := controllers.FetchOrCreateUser(email, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to sync user",
		})
	}

	claims := fiber.Map{
		"sub":      introspection.Sub,
		"email":    email,
		"name":     name,
		"scope":    introspection.Scope,
		"clientID": introspection.ClientID,
		"exp":      introspection.Exp,
	}

	c.Locals("claims", claims)
	c.Locals("userID", userID)
	c.Locals("kratosID", introspection.Sub)

	return c.Next()
}
