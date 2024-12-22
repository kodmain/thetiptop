package repositories

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	errors_domain_game "github.com/kodmain/thetiptop/api/internal/domain/game/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

type GameRepository struct {
	store *database.Database
}

type GameRepositoryInterface interface {
	// Ticket
	CreateTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface)
	CreateTickets(objs []*transfert.Ticket, options ...database.Option) errors.ErrorInterface
	ReadTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface)
	ReadTickets(obj *transfert.Ticket, options ...database.Option) ([]*entities.Ticket, errors.ErrorInterface)
	UpdateTicket(entity *entities.Ticket, options ...database.Option) errors.ErrorInterface
	DeleteTicket(obj *transfert.Ticket, options ...database.Option) errors.ErrorInterface
	CountTicket(obj *transfert.Ticket, options ...database.Option) (int, errors.ErrorInterface)
}

func NewGameRepository(store *database.Database) *GameRepository {
	store.Engine.AutoMigrate(entities.Ticket{})
	return &GameRepository{store}
}

// CreateTicket creates a new ticket
// Inserts a new ticket into the database based on the transfert.Ticket input object
//
// Parameters:
// - obj: *transfert.Ticket - The ticket transfer object to create
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - *entities.Ticket: The created ticket entity
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) CreateTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface) {
	ticket := entities.CreateTicket(obj)

	// Applique les options à la requête
	query := r.store.Engine.Create(ticket)
	for _, option := range options {
		option(query)
	}

	if query.Error != nil {
		return nil, errors.ErrInternalServer.Log(query.Error)
	}

	return ticket, nil
}

// CreateTickets creates multiple tickets
// Inserts multiple tickets into the database in a single batch operation
//
// Parameters:
// - objs: []*transfert.Ticket - The slice of ticket transfer objects to create
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - []*entities.Ticket: The slice of created ticket entities
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) CreateTickets(objs []*transfert.Ticket, options ...database.Option) errors.ErrorInterface {
	tickets := make([]*entities.Ticket, len(objs))
	for i, obj := range objs {
		tickets[i] = entities.CreateTicket(obj)
	}

	query := r.store.Engine.CreateInBatches(tickets, len(tickets))
	for _, option := range options {
		option(query)
	}

	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

// ReadTickets reads multiple tickets from the database
// Finds and returns a list of tickets based on the provided transfer object and options
//
// Parameters:
// - obj: *transfert.Ticket - The ticket transfer object with search parameters
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - []*entities.Ticket: A slice of found ticket entities
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) ReadTickets(obj *transfert.Ticket, options ...database.Option) ([]*entities.Ticket, errors.ErrorInterface) {
	var tickets []*entities.Ticket

	query := r.store.Engine.Where(obj)
	for _, option := range options {
		option(query)
	}

	result := query.Find(&tickets)

	if result.Error != nil {
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return tickets, nil
}

// ReadTicket reads a ticket from the database
// Finds and returns a ticket based on the provided transfer object and options
//
// Parameters:
// - obj: *transfert.Ticket - The ticket transfer object with search parameters
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - *entities.Ticket: The found ticket entity
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) ReadTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface) {
	ticket := &entities.Ticket{}

	query := r.store.Engine.Where(obj)
	for _, option := range options {
		option(query)
	}

	result := query.First(ticket)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors_domain_game.ErrTicketNotFound
		}
		return nil, errors.ErrInternalServer.Log(result.Error)
	}

	return ticket, nil
}

// UpdateTicket updates an existing ticket in the database
// Saves the updated ticket entity
//
// Parameters:
// - entity: *entities.Ticket - The ticket entity to update
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) UpdateTicket(entity *entities.Ticket, options ...database.Option) errors.ErrorInterface {
	query := r.store.Engine.Save(entity)
	for _, option := range options {
		option(query)
	}

	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

// DeleteTicket deletes a ticket from the database
// Removes a ticket based on the provided transfer object
//
// Parameters:
// - obj: *transfert.Ticket - The ticket transfer object to delete
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) DeleteTicket(obj *transfert.Ticket, options ...database.Option) errors.ErrorInterface {
	ticket := entities.CreateTicket(obj)
	query := r.store.Engine.Where(obj).Delete(ticket)
	for _, option := range options {
		option(query)
	}

	if query.Error != nil {
		return errors.ErrInternalServer.Log(query.Error)
	}

	return nil
}

// CountTicket counts the number of tickets in the database
// Returns the number of tickets based on the provided transfer object and options
//
// Parameters:
// - obj: *transfert.Ticket - The ticket transfer object with search parameters
// - options: ...database.Option - Additional options to customize the query
//
// Returns:
// - int: The number of tickets found
// - errors.ErrorInterface: The error interface if an error occurs
func (r *GameRepository) CountTicket(obj *transfert.Ticket, options ...database.Option) (int, errors.ErrorInterface) {
	ticket := entities.CreateTicket(obj)
	var count int64

	query := r.store.Engine.Model(&ticket).Where(ticket)
	for _, option := range options {
		option(query)
	}

	result := query.Count(&count)

	if result.Error != nil {
		return 0, errors.ErrInternalServer.Log(result.Error)
	}

	return int(count), nil
}
