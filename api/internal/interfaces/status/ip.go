package status

import "github.com/gofiber/fiber/v2"

// @Summary		Show the ip of user.
// @Description	get the ip of user.
// @Tags		Status
// @Accept		*/*
// @Produce		application/json
// @Success		200	{object}	nil
// @Router		/status/ip [get]
// @Id	        status.IP
func IP(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString(c.IP())
}
