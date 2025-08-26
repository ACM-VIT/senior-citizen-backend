package middleware

import (
	"acm-portal/helpers"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// extractToken gets the JWT from Authorization header or cookie
func extractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return c.Cookies(os.Getenv("CSRF_COOKIE_NAME"))
}

// AuthMiddleware validates the JWT and sets user info in c.Locals
func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := extractToken(c)
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	// Use helper to parse claims
	claims, err := helpers.ParseJWTClaims(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Invalid token",
			"details": err.Error(),
		})
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token expired",
		})
	}

	// Attach user info to request context
	c.Locals("userID", claims["ID"])
	c.Locals("email", claims["email"])
	c.Locals("admin", claims["admin"])

	return c.Next()
}
