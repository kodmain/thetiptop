package services

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
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
	GetTickets() ([]*entities.Ticket, errors.ErrorInterface)
	GetRandomTicket() (*entities.Ticket, errors.ErrorInterface)
	UpdateTicket(*transfert.Ticket) (*entities.Ticket, errors.ErrorInterface)
}
