package jwt

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	token := c.Locals("token").(*Token)
	if token == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "No token")
	}

	if token.HasExpired() {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return c.Next()
}

func Parser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header format")
	}

	tokenString := parts[1]
	token, err := tokenToClaims(tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token: "+err.Error())
	}

	c.Locals("token", token)

	return c.Next()
}
