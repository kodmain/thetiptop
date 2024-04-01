package client

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/repositories"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
)

// @Tags		Sign
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		email		formData	string	true	"Email address"
// @Success		204	{object}	nil
// @Router		/sign/up [post]
// @Id			client.SignUp
func SignUp(c *fiber.Ctx) error {
	email := c.FormValue("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}

	service := services.NewClientService(
		repositories.NewClientRepository(),
	)

	service.CreateClient(&entities.Client{
		Email: email,
	})

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Param		email		formData	string	true	"Email address"
// @Success		204	{object}	nil
// @Router		/sign/in [post]
// @Id			client.SignIn
func SignIn(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Sign
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/sign/out [get]
// @Id			client.SignOut
func SignOut(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
