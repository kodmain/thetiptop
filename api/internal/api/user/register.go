package user

import "github.com/gofiber/fiber/v2"

// @Summary		Show the user register action.
// @Description	Register a new user in the system.
// @Tags		User
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/user/register [post]
// @Id	        user.Register
func Register(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
