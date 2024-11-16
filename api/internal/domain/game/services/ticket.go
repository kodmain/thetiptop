package services

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

func (s *GameService) GetRandomTicket() (*entities.Ticket, errors.ErrorInterface) {
	ticket, err := s.repo.ReadTicket(&transfert.Ticket{}, database.Where("client_id IS NULL"), database.Order("RANDOM()"))
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *GameService) Validate() (bool, errors.ErrorInterface) {
	return true, nil
}
