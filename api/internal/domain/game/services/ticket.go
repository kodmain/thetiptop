package services

import (
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func (s *GameService) Lucky() (*entities.Ticket, errors.ErrorInterface) {
	return &entities.Ticket{}, nil
}

func (s *GameService) Validate() (bool, errors.ErrorInterface) {
	return true, nil
}
