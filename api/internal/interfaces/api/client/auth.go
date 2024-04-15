package client

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/services"

	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

// @Tags		Sign
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address"
// @Success		201	{object}	nil "Client created"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		409	{object}	nil "Client already exists"
// @Router		/sign/up [post]
// @Id			client.SignUp
func SignUp(c *fiber.Ctx) error {
	status, response := services.SignUp(
		c.FormValue("email"),
		c.FormValue("password"),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Param		email		formData	string	true	"Email address"
// @Param		password	formData	string	true	"Password"
// @Success		204	{object}	nil
// @Router		/sign/in [post]
// @Id			client.SignIn
func SignIn(c *fiber.Ctx) error {
	status, response := services.SignIn(
		c.FormValue("email"),
		c.FormValue("password"),
	)

	return c.Status(status).JSON(response)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/sign/renew [get]
// @Id			client.SignRenew
func SignRenew(c *fiber.Ctx) error {
	token := c.Locals("token")
	if token == nil {
		return fiber.NewError(fiber.StatusBadRequest, "No token")
	}

	status, response := services.SignRenew(token.(*serializer.Token))

	return c.Status(status).JSON(response)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/sign/out [get]
// @Id			jwt.Auth => client.SignOut
func SignOut(c *fiber.Ctx) error {
	// TODO: Implement SignOut
	return c.Status(fiber.StatusNoContent).Send(nil)
}
