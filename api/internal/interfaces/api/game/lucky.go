package game

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/application/services/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/repositories"
	"github.com/kodmain/thetiptop/api/internal/domain/game/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

// @Tags		Game
// @Accept		multipart/form-data
// @Summary		List all code errors.
// @Produce		application/json
// @Router		/game/lucky [get]
// @Id			game.Lucky
func Lucky(ctx *fiber.Ctx) error {
	status, response := game.Lucky(
		services.Game(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewGameRepository(database.Get(config.GetString("services.game.database", config.DEFAULT))),
		),
	)

	return ctx.Status(status).JSON(response)
}
