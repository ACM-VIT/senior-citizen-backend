package middleware

import "github.com/gofiber/fiber/v2"

func AdminMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        return c.Next()
    }
}
