package client

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/workflow"
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
	password := c.FormValue("password")

	status, response := workflow.Client.SignUp(email, password)

	return c.Status(status).JSON(response)
	/*
		if err := workflow.SignUp(email, password); logger.Error(err) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusNoContent).Send(nil)
	*/
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
