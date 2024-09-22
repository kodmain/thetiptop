package user

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"

	"github.com/kodmain/thetiptop/api/internal/domain/user/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
)

// @Tags		Client
// @Accept		multipart/form-data
// @Summary		Register a client.
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Param		password	formData	string	true	"Password" default(Aa1@azetyuiop)
// @Param 		cgu			formData	bool	true	"CGU" default(true)
// @Param 		newsletter	formData	bool	true	"Newsletter" default(false)
// @Success		201	{object}	nil "Client created"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		409	{object}	nil "Client already exists"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/client/register [post]
// @Id			user.RegisterClient
func RegisterClient(ctx *fiber.Ctx) error {
	dtoCredential := &transfert.Credential{}
	if err := ctx.BodyParser(dtoCredential); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dtoClient := &transfert.Client{}
	if err := ctx.BodyParser(dtoClient); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.RegisterClient(
		domain.User(
			ctx,
			repositories.NewUserRepository(database.Get(config.GetString("services.client.database", config.DEFAULT))),
			mail.Get(config.GetString("services.client.mail", config.DEFAULT)),
		), dtoCredential, dtoClient,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags		Client
// @Accept		multipart/form-data
// @Summary		Update a client.
// @Produce		application/json
// @Param		id			formData	string	true	"Client ID" format(uuid)
// @Param		newsletter	formData	bool	true	"Newsletter" default(false)
// @Success		204	{object}	nil "Password updated"
// @Failure		400	{object}	nil "Invalid email, password or token"
// @Failure		404	{object}	nil "Client not found"
// @Failure		409	{object}	nil "Client already validated"
// @Failure		410	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/client [put]
// @Id			jwt.Auth => user.UpdateClient
func UpdateClient(ctx *fiber.Ctx) error {
	dtoClient := &transfert.Client{}
	if err := ctx.BodyParser(dtoClient); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.UpdateClient(
		domain.User(
			ctx,
			repositories.NewUserRepository(database.Get(config.GetString("services.client.database", config.DEFAULT))),
			mail.Get(config.GetString("services.client.mail", config.DEFAULT)),
		), dtoClient,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags		Client
// @Accept		multipart/form-data
// @Summary		Get a client by ID.
// @Produce		application/json
// @Param		id			path		string	true	"Client ID" format(uuid)
// @Success		200	{object}	nil "Client details"
// @Failure		400	{object}	nil "Invalid client ID"
// @Failure		404	{object}	nil "Client not found"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/client/{id} [get]
// @Id			jwt.Auth => user.GetClient
func GetClient(ctx *fiber.Ctx) error {
	clientID := ctx.Params("id")

	if clientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Client ID is required")
	}

	dtoClient := &transfert.Client{
		ID: &clientID,
	}

	status, response := services.GetClient(
		domain.User(
			ctx,
			repositories.NewUserRepository(database.Get(config.GetString("services.client.database", config.DEFAULT))),
			mail.Get(config.GetString("services.client.mail", config.DEFAULT)),
		), dtoClient,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags		Client
// @Summary		Delete a client by ID.
// @Produce		application/json
// @Param		id			path		string	true	"Client ID" format(uuid)
// @Success		204	{object}	nil "Client deleted"
// @Failure		400	{object}	nil "Invalid client ID"
// @Failure		404	{object}	nil "Client not found"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/client/{id} [delete]
// @Id			jwt.Auth => user.DeleteClient
func DeleteClient(ctx *fiber.Ctx) error {
	clientID := ctx.Params("id")

	if clientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Client ID is required")
	}

	dtoClient := &transfert.Client{
		ID: &clientID,
	}

	status, response := services.DeleteClient(
		domain.User(
			ctx,
			repositories.NewUserRepository(database.Get(config.GetString("services.client.database", config.DEFAULT))),
			mail.Get(config.GetString("services.client.mail", config.DEFAULT)),
		), dtoClient,
	)

	return ctx.Status(status).JSON(response)
}
