package jwt

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func Auth(c *fiber.Ctx) error {
	auth := c.Locals("token")
	if auth == nil {
		return c.Status(errors.ErrAuthNoToken.Code()).JSON(errors.ErrAuthNoToken)
	}

	token := auth.(*Token)
	if token.HasExpired() {
		return c.Status(errors.ErrAuthExpiredToken.Code()).JSON(errors.ErrAuthExpiredToken)
	}

	if token.IsNotValid() {
		return c.Status(errors.ErrAuthInvalidToken.Code()).JSON(errors.ErrAuthInvalidToken)
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
		return c.Status(errors.ErrAuthBadFormat.Code()).JSON(errors.ErrAuthBadFormat)
	}

	tokenString := parts[1]
	token, err := TokenToClaims(tokenString)
	if err != nil {
		return c.Status(errors.ErrAuthFailed.Code()).JSON(errors.ErrAuthFailed)
	}

	c.Locals("token", token)

	return c.Next()
}
