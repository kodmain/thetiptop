package client

import (
	"github.com/gofiber/fiber/v2"
)

// @Tags		Models
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client/:id [get]
// @Id			client.FindOne
func FindOne(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Models
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client [get]
// @Id			client.Find
func Find(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Models
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client [patch]
// @Id			client.UpdatePartial
func UpdatePartial(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Models
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client [put]
// @Id			client.UpdateComplete
func UpdateComplete(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// @Tags		Models
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/client [delete]
// @Id			client.Delete
func Delete(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
