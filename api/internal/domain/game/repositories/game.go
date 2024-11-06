package repositories

import (
	"sync"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

var game sync.Once

type UserRepository struct {
}
type GameRepository struct {
	store *database.Database
}

type GameRepositoryInterface interface {
	// Ticket
	CreateTicket(obj *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface)
	ReadTicket(obj *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface)
	UpdateTicket(entity *entities.Ticket) errors.ErrorInterface
	DeleteTicket(obj *transfert.Ticket) errors.ErrorInterface
}

func NewGameRepository(store *database.Database) *GameRepository {
	game.Do(func() {
		store.Engine.AutoMigrate(entities.Ticket{})
	})

	return &GameRepository{store}
}

func (r *GameRepository) CreateTicket(obj *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface) {

	return nil, nil
}

func (r *GameRepository) ReadTicket(obj *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface) {

	return nil, nil
}

func (r *GameRepository) UpdateTicket(entity *entities.Ticket) errors.ErrorInterface {
	return nil
}

func (r *GameRepository) DeleteTicket(obj *transfert.Ticket) errors.ErrorInterface {
	return nil
}
