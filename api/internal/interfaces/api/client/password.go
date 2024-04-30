package client

import "github.com/gofiber/fiber/v2"

// @Tags		Password
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client/password/renew [post]
// @Id			client.Renew
func Renew(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Password
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client/password/reset [post]
// @Id			client.Reset
func Reset(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
