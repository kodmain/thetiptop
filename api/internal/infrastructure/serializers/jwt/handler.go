package jwt

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func Auth(c *fiber.Ctx) error {
	auth := c.Locals("token")
	if auth == nil {
		return fiber.NewError(errors.ErrAuthNoToken.Code(), errors.ErrAuthNoToken.Error())
	}

	token := auth.(*Token)
	if token.HasExpired() {
		return fiber.NewError(errors.ErrAuthExpiredToken.Code(), errors.ErrAuthExpiredToken.Error())
	}

	if token.IsNotValid() {
		return fiber.NewError(errors.ErrAuthInvalidToken.Code(), errors.ErrAuthInvalidToken.Error())
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
		return fiber.NewError(errors.ErrAuthBadFormat.Code(), errors.ErrAuthBadFormat.Error())
	}

	tokenString := parts[1]
	token, err := TokenToClaims(tokenString)
	if err != nil {
		return fiber.NewError(errors.ErrAuthFailed.Code(), errors.ErrAuthFailed.Error())
	}

	c.Locals("token", token)

	return c.Next()
}
