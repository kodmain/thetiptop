package game

import (
	"github.com/gofiber/fiber/v2"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/services"
)

func GetRandomTicket(service services.GameServiceInterface) (int, any) {
	ticket, err := service.GetRandomTicket()
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, ticket
}

func GetTickets(service services.GameServiceInterface) (int, any) {
	tickets, err := service.GetTickets()
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, tickets
}

func UpdateTicket(service services.GameServiceInterface, dtoTicket *transfert.Ticket) (int, any) {
	ticket, err := service.UpdateTicket(dtoTicket)

	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, ticket
}

func GetTicketById(service services.GameServiceInterface, dtoTicket *transfert.Ticket) (int, any) {
	ticket, err := service.GetTicketById(dtoTicket)

	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, ticket
}
