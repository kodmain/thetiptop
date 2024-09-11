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
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

// @Tags		Client
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Param		password	formData	string	true	"Password" default(Aa1@azetyuiop)
// @Param 		cgu			formData	bool	true	"CGU" default(true)
// @Param 		newsletter	formData	bool	true	"Newsletter" default(false)
// @Success		201	{object}	nil "Client created"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		409	{object}	nil "Client already exists"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/user/register [post]
// @Id			client.Register
func Register(c *fiber.Ctx) error {
	dtoCredential := &transfert.Credential{}
	if err := c.BodyParser(dtoCredential); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dtoClient := &transfert.Client{}
	if err := c.BodyParser(dtoClient); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.UserRegister(
		domain.Client(
			repositories.NewClientRepository(database.Get(config.Get("services.client.database").(string))),
			mail.Get(),
		), dtoCredential, dtoClient,
	)

	return c.Status(status).JSON(response)
}

// @Tags		Client
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Param		password	formData	string	true	"Password" default(Aa1@azetyuiop)
// @Success		200	{object}	nil "Client signed in"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/user/auth [post]
// @Id			client.UserAuth
func UserAuth(c *fiber.Ctx) error {
	dto := &transfert.Credential{}
	if err := c.BodyParser(dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.UserAuth(
		domain.Client(
			repositories.NewClientRepository(database.Get(config.Get("services.client.database").(string))),
			mail.Get(),
		), dto,
	)

	return c.Status(status).JSON(response)
}

// @Tags		Client
// @Accept		*/*
// @Accept		multipart/form-data
// @Produce		application/json
// @Success		200	{object}	nil "JWT token renewed"
// @Failure		400	{object}	nil "Invalid token"
// @Failure		401	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Param 		Authorization header string true "With the bearer started"
// @Router		/user/auth/renew [get]
// @Id			client.UserAuthRenew
func UserAuthRenew(c *fiber.Ctx) error {
	token := c.Locals("token")
	if token == nil {
		return fiber.NewError(fiber.StatusBadRequest, "no token")
	}

	status, response := services.UserAuthRenew(
		token.(*jwt.Token),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Client
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Param		password	formData	string	true	"Password" default(Aa1@azetyuiop)
// @Param		token		formData	string	true	"Token"
// @Success		204	{object}	nil "Password updated"
// @Failure		400	{object}	nil "Invalid email, password or token"
// @Failure		404	{object}	nil "Client not found"
// @Failure		409	{object}	nil "Client already validated"
// @Failure		410	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/user/password [put]
// @Id			client.PasswordUpdate
func PasswordUpdate(c *fiber.Ctx) error {
	dtoCredential := &transfert.Credential{}
	if err := c.BodyParser(dtoCredential); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dtoValidation := &transfert.Validation{}
	if err := c.BodyParser(dtoValidation); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.CredentialUpdate(
		domain.Client(
			repositories.NewClientRepository(database.Get(config.Get("services.client.database").(string))),
			mail.Get(),
		), dtoValidation, dtoCredential,
	)

	return c.Status(status).JSON(response)
}
