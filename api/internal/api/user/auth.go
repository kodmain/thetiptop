package user

import "github.com/gofiber/fiber/v2"

// @Summary		Show the user auth action.
// @Description	Auth user in the system.
// @Tags		User
// @Accept		*/*
// @Produce		application/json
// @Success		200	{object}	model.User
// @Router		/user/auth [post]
// @Id	        user.Auth
func Auth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
