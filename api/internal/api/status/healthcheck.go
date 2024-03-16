// Package status implementation all handler for Status API
package status

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck is an HTTP handler function that returns a "No Content" response to indicate that the
// API is running and healthy. It is typically used as a "health check" endpoint that can be polled by
// monitoring systems to verify that the API is functioning correctly.
// Parameters:
// - c: a pointer to the fiber.Ctx object representing the HTTP request context.
// Returns:
// - an error value, which is always nil, since there is no meaningful error condition for this endpoint.
// @Summary		Show the status of server.
// @Description	get the status of server.
// @Tags		Status
// @Accept		*/*
// @Produce		application/json
// @Success		204	{object}	nil
// @Router		/status/healthcheck [get]
// @Id	        status.HealthCheck
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}
