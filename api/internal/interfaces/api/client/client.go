package client

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodmain/thetiptop/api/internal/application/services"

	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

// @Tags		Sign
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email)
// @Param		password	formData	string	true	"Password"
// @Success		201	{object}	nil "Client created"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		409	{object}	nil "Client already exists"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/sign/up [post]
// @Id			client.SignUp
func SignUp(c *fiber.Ctx) error {
	status, response := services.SignUp(
		domain.Client(
			repositories.NewClientRepository(database.Get()),
			mail.Get(),
		),
		c.FormValue("email"),
		c.FormValue("password"),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Param		token	formData	string	true	"Token"
// @Param		email	formData	string	true	"Email address" format(email)
// @Success		204	{object}	nil "Client email validate"
// @Failure		400	{object}	nil "Invalid email or token"
// @Failure		404	{object}	nil "Client not found"
// @Failure		409	{object}	nil "Client already validated"
// @Failure		410 {object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/sign/validation [put]
// @Id			client.SignValidation
func SignValidation(c *fiber.Ctx) error {
	status, response := services.SignValidation(
		domain.Client(
			repositories.NewClientRepository(database.Get()),
			mail.Get(),
		),
		c.FormValue("email"),
		c.FormValue("token"),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email)
// @Param		password	formData	string	true	"Password"
// @Success		200	{object}	nil "Client signed in"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/sign/in [post]
// @Id			client.SignIn
func SignIn(c *fiber.Ctx) error {
	status, response := services.SignIn(
		domain.Client(
			repositories.NewClientRepository(database.Get()),
			mail.Get(),
		),
		c.FormValue("email"),
		c.FormValue("password"),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Success		200	{object}	nil "JWT token renewed"
// @Failure		400	{object}	nil "Invalid token"
// @Failure		401	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/sign/renew [get]
// @Id			client.SignRenew
func SignRenew(c *fiber.Ctx) error {
	token := c.Locals("token")
	if token == nil {
		return fiber.NewError(fiber.StatusBadRequest, "no token")
	}

	status, response := services.SignRenew(
		token.(*serializer.Token),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Password
// @Accept		*/*
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email)
// @Success		204	{object}	nil "Password recover"
// @Failure		400	{object}	nil "Invalid email"
// @Failure		404	{object}	nil "Client not found"
// @Failure		500	{object}	nil "Internal server error"
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

// @Tags		Password
// @Accept		*/*
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email)
// @Param		password	formData	string	true	"Password"
// @Param		token		formData	string	true	"Token"
// @Success		204	{object}	nil "Password updated"
// @Failure		400	{object}	nil "Invalid email, password or token"
// @Failure		404	{object}	nil "Client not found"
// @Failure		409	{object}	nil "Client already validated"
// @Failure		410	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/password/update [put]
// @Id			client.PasswordUpdate
func PasswordUpdate(c *fiber.Ctx) error {
	status, response := services.PasswordUpdate(
		domain.Client(
			repositories.NewClientRepository(database.Get()),
			mail.Get(),
		),
		c.FormValue("email"),
		c.FormValue("password"),
		c.FormValue("token"),
	)

	return c.Status(status).JSON(response)
}
