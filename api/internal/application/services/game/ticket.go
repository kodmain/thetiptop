package game

import (
	"github.com/gofiber/fiber/v2"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/services"
)

func Lucky(service services.GameServiceInterface) (int, any) {
	ticket, err := service.Lucky()
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, ticket
}

func Validate(service services.GameServiceInterface, dtoTicket *transfert.Ticket) (int, any) {
	return fiber.StatusNoContent, nil
}
