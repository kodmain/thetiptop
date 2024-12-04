package game_test

import (
	"sync"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/mock"
)

// DomainGameService is a mock implementation of the GameServiceInterface
// This mock is used to simulate the behavior of the service for testing purposes.
type DomainGameService struct {
	mock.Mock
	mu sync.Mutex // Ensures thread safety when the mock is used in concurrent tests.
}

// GetRandomTicket simulates the GetRandomTicket method of the GameServiceInterface
//
// It uses testify's mock functionality to simulate return values and errors.
//
// Returns:
// - *entities.Ticket: the ticket returned by the service, if successful
// - errors.ErrorInterface: the error returned by the service, if any
func (mgs *DomainGameService) GetRandomTicket() (*entities.Ticket, errors.ErrorInterface) {
	args := mgs.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Ticket), nil
}

// UpdateTicket simulates the UpdateTicket method of the GameServiceInterface
//
// It uses testify's mock functionality to simulate return values and errors.
//
// Parameters:
// - dtoTicket: *game.Ticket - the ticket to be updated
//
// Returns:
// - *entities.Ticket: the updated ticket, if successful
// - errors.ErrorInterface: the error returned by the service, if any
func (mgs *DomainGameService) UpdateTicket(dtoTicket *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface) {
	args := mgs.Called(dtoTicket)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Ticket), nil
}

// GetTickets simulates the GetTickets method of the GameServiceInterface
//
// It uses testify's mock functionality to simulate return values and errors.
//
// Parameters:
// - dtoTicket: *game.Ticket - the ticket to be updated
//
// Returns:
// - *entities.Ticket: the updated ticket, if successful
// - errors.ErrorInterface: the error returned by the service, if any
func (mgs *DomainGameService) GetTickets() ([]*entities.Ticket, errors.ErrorInterface) {
	args := mgs.Called()
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).([]*entities.Ticket), nil
}

// GetTicketById simulates the GetTicketById method of the GameServiceInterface
//
// It uses testify's mock functionality to simulate return values and errors.
//
// Parameters:
// - dtoTicket: *game.Ticket - the ticket to be updated
//
// Returns:
// - *entities.Ticket: the updated ticket, if successful
// - errors.ErrorInterface: the error returned by the service, if any
func (mgs *DomainGameService) GetTicketById(dtoTicket *transfert.Ticket) (*entities.Ticket, errors.ErrorInterface) {
	args := mgs.Called(dtoTicket)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Ticket), nil
}
