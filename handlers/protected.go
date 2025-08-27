package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
)

// Middleware: validate access_token
func requireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid Authorization format"})
	}

	tokenString := parts[1]

	// ⚠️ TODO: verify using Hydra JWKS
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)

	// Store claims in context for later use
	c.Locals("claims", claims)

	return c.Next()
}

// Protected route(s)
func ProtectedRouter(app *fiber.App) {
	app.Get("/protected", requireAuth, func(c *fiber.Ctx) error {
		claims := c.Locals("claims")
		return c.JSON(fiber.Map{
			"message": "You are authenticated 🎉",
			"claims":  claims,
		})
	})
}