package game

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/application/services/game"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/repositories"
	"github.com/kodmain/thetiptop/api/internal/domain/game/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

// @Tags		Game
// @Accept		multipart/form-data
// @Summary		Get a lucky ticket.
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

// @Tags		Game
// @Accept		multipart/form-data
// @Summary	  	Reveal the prize of a ticket.
// @Produce		application/json
// @Router		/game/reveal [post]
// @Id			game.Reveal
func Reveal(ctx *fiber.Ctx) error {
	dtoTicket := &transfert.Ticket{}
	if err := ctx.BodyParser(dtoTicket); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := game.Validate(
		services.Game(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewGameRepository(database.Get(config.GetString("services.game.database", config.DEFAULT))),
		), dtoTicket,
	)

	return ctx.Status(status).JSON(response)
}
