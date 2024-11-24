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
// @Summary		Get a random ticket.
// @Produce		application/json
// @Router		/game/ticket/random [get]
// @Id			jwt.Auth => game.GetTicket
// @Security 	Bearer
// @Success		200	{object} 	nil "Ticket details"
// @Failure		400	{object} 	nil "Bad request"
// @Failure		401	{object} 	nil "Unauthorized"
func GetTicket(ctx *fiber.Ctx) error {
	status, response := game.GetRandomTicket(
		services.Game(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewGameRepository(database.Get(config.GetString("services.game.database", config.DEFAULT))),
		),
	)

	return ctx.Status(status).JSON(response)
}

// @Tags		Game
// @Accept		multipart/form-data
// @Summary		Get a random ticket.
// @Produce		application/json
// @Router		/game/ticket [get]
// @Id			jwt.Auth => game.GetTickets
// @Security 	Bearer
// @Success		200	{object} 	nil "Tickets details"
// @Failure		400	{object} 	nil "Bad request"
// @Failure		404	{object} 	nil "Not found"
func GetTickets(ctx *fiber.Ctx) error {
	status, response := game.GetTickets(
		services.Game(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewGameRepository(database.Get(config.GetString("services.game.database", config.DEFAULT))),
		),
	)

	return ctx.Status(status).JSON(response)
}

// @Tags		Game
// @Accept		multipart/form-data
// @Summary	  	Update a ticket.
// @Produce		application/json
// @Router		/game/ticket [put]
// @Id			jwt.Auth => game.UpdateTicket
// @Security 	Bearer
// @Param		id	formData	string	true	"Ticket ID" format(uuid)
// @Success		200	{object} 	nil "Ticket details"
// @Failure		400	{object} 	nil "Bad request"
// @Failure		401	{object} 	nil "Unauthorized"
// @Failure		404	{object} 	nil "Not found"
func UpdateTicket(ctx *fiber.Ctx) error {
	dtoTicket := &transfert.Ticket{}
	if err := ctx.BodyParser(dtoTicket); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := game.UpdateTicket(
		services.Game(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewGameRepository(database.Get(config.GetString("services.game.database", config.DEFAULT))),
		), dtoTicket,
	)

	return ctx.Status(status).JSON(response)
}
