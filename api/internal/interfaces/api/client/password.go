package client

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
)

// @Tags		Password
// @Accept		*/*
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email)
// @Success		201	{object}	nil
// @Failure		400	{object}	nil
// @Router		/password/recover [post]
// @Id			client.PasswordRecover
func PasswordRecover(c *fiber.Ctx) error {
	status, response := services.PasswordRecover(
		domain.Client(
			repositories.NewClientRepository(database.Get()),
			mail.Get(),
		),
		c.FormValue("email"),
	)

	return c.Status(status).JSON(response)
}
