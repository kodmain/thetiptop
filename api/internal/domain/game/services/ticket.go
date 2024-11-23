package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	user "github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

func (s *GameService) GetRandomTicket() (*entities.Ticket, errors.ErrorInterface) {
	if !s.security.IsGrantedByRoles(user.ROLE_EMPLOYEE) {
		return nil, errors.ErrUnauthorized
	}

	ticket, err := s.repo.ReadTicket(&transfert.Ticket{}, database.Where("client_id IS NULL"), database.Order("RANDOM()"))
	if err != nil {
		return nil, errors.ErrNoData
	}

	return ticket, nil
}

func (s *GameService) GetTickets() ([]*entities.Ticket, errors.ErrorInterface) {
	if !s.security.IsGrantedByRules() {
		return nil, errors.ErrUnauthorized
	}

	tickets, err := s.repo.ReadTickets(&transfert.Ticket{
		ClientID: s.security.GetCredentialID(),
	})

	if err != nil {
		return nil, errors.ErrNoData
	}

	return tickets, nil
}

func (s *GameService) UpdateTicket(dto *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface) {
	ticket, err := s.repo.ReadTicket(dto, database.Where("client_id IS NULL"))

	if err != nil {
		return nil, err
	}

	if !s.security.CanUpdate(ticket) {
		return nil, errors.ErrUnauthorized
	}

	ticket.ClientID = dto.ClientID

	if err := s.repo.UpdateTicket(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
