package middleware

import (
	"acm-portal/controllers"
	"encoding/json"
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
	Name string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Admin    bool     `json:"admin,omitempty"`
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

    if !introspection.Active || introspection.Exp < time.Now().Unix() {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Inactive or expired token",
        })
    }

    userID, err := controllers.FetchOrCreateUser(introspection.Email, introspection.Name)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to sync user",
        })
    }

	isAdmin, err := controllers.FetchOrCreateAdmin(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to check admin status",
		})
	}

    claims := fiber.Map{
        "sub":      introspection.Sub,
        "email":    introspection.Email,
        "name":     introspection.Name,
        "admin":    introspection.Admin,
        "scope":    introspection.Scope,
        "clientID": introspection.ClientID,
        "exp":      introspection.Exp,
    }

    c.Locals("claims", claims)
    c.Locals("userID", userID) 
    c.Locals("kratosID", introspection.Sub)
	c.Locals("admin", isAdmin)

    return c.Next()
}
