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
// @Param		token	formData	string	true	"Token"
// @Param		email	formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Success		204	{object}	nil "Client email validate"
// @Failure		400	{object}	nil "Invalid email or token"
// @Failure		404	{object}	nil "Client not found"
// @Failure		409	{object}	nil "Client already validated"
// @Failure		410 {object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/user/register/validation [put]
// @Id			client.SignValidation
func SignValidation(c *fiber.Ctx) error {
	dtoCredential := &transfert.Credential{}
	if err := c.BodyParser(dtoCredential); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dtoValidation := &transfert.Validation{}
	if err := c.BodyParser(dtoValidation); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.SignValidation(
		domain.Client(
			repositories.NewClientRepository(database.Get(config.Get("services.client.database").(string))),
			mail.Get(),
		), dtoValidation, dtoCredential,
	)

	return c.Status(status).JSON(response)
}

// @Tags		Client
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Param		type		formData	string	true	"Type of validation" enums(mail, password, phone)
// @Router		/validation/renew [post]
// @Id			client.ValidationRecover
func ValidationRecover(c *fiber.Ctx) error {
	dtoCredential := &transfert.Credential{}
	if err := c.BodyParser(dtoCredential); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dtoValidation := &transfert.Validation{}
	if err := c.BodyParser(dtoValidation); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.ValidationRecover(
		domain.Client(
			repositories.NewClientRepository(database.Get(config.Get("services.client.database").(string))),
			mail.Get(),
		), dtoCredential, dtoValidation,
	)

	return c.Status(status).JSON(response)
}
