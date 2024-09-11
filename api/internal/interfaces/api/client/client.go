package client

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"

	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
)

// @Tags		Client
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		id			formData	string	true	"Client ID" format(uuid)
// @Param		newsletter	formData	bool	true	"Newsletter" default(false)
// @Success		204	{object}	nil "Password updated"
// @Failure		400	{object}	nil "Invalid email, password or token"
// @Failure		404	{object}	nil "Client not found"
// @Failure		409	{object}	nil "Client already validated"
// @Failure		410	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/user/update [put]
// @Id			client.UpdateClient
func UpdateClient(c *fiber.Ctx) error {
	dtoClient := &transfert.Client{}
	if err := c.BodyParser(dtoClient); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.UpdateClient(
		domain.Client(
			repositories.NewClientRepository(database.Get(config.Get("services.client.database").(string))),
			mail.Get(),
		), dtoClient,
	)

	return c.Status(status).JSON(response)
}
