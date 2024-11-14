package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *GameService) Lucky() (*entities.Ticket, errors.ErrorInterface) {
	ticket, err := s.repo.ReadTicket(&transfert.Ticket{})
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *GameService) Validate() (bool, errors.ErrorInterface) {
	return true, nil
}
