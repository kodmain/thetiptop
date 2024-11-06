package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/game/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type GameService struct {
	security security.PermissionInterface
	repo     repositories.GameRepositoryInterface
}

func Game(security security.PermissionInterface, repo repositories.GameRepositoryInterface) *GameService {
	return &GameService{security, repo}
}

type GameServiceInterface interface {
	Lucky() (*entities.Ticket, errors.ErrorInterface)
	Validate() (bool, errors.ErrorInterface)
}
