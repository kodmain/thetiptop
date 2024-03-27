// Package status implementation all handler for Status API
package status

import (
	"github.com/gofiber/fiber/v2"
)

// @Summary		Show the status of server.
// @Description	get the status of server.
// @Tags		Status
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/status/healthcheck [get]
// @Id	        metrics.Counter => status.HealthCheck
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
