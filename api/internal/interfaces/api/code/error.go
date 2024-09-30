package code

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/domain/code/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/code/services"
)

// @Tags		Error
// @Accept		multipart/form-data
// @Summary		List all code errors.
// @Produce		application/json
// @Router		/code/error [get]
// @Id			code.ListErrors
func ListErrors(ctx *fiber.Ctx) error {
	status, response := services.ListErrors(
		domain.Code(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewCodeRepository(),
		),
	)

	return ctx.Status(status).JSON(response)
}
