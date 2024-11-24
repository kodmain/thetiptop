package events_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/game/events"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

// MockGameRepository simule l'interface GameRepositoryInterface
type MockGameRepository struct {
	mock.Mock
}

// CreateTicket simule la création d'un ticket
func (m *MockGameRepository) CreateTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Error(1) != nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Ticket), nil
}

// CreateTickets simule la création de plusieurs tickets
func (m *MockGameRepository) CreateTickets(objs []*transfert.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)
	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// ReadTicket simule la lecture d'un ticket
func (m *MockGameRepository) ReadTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Error(1) != nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Ticket), nil
}

// ReadTickets simule la lecture de plusieurs tickets
func (m *MockGameRepository) ReadTickets(obj *transfert.Ticket, options ...database.Option) ([]*entities.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Error(1) != nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).([]*entities.Ticket), nil
}

// UpdateTicket simule la mise à jour d'un ticket
func (m *MockGameRepository) UpdateTicket(entity *entities.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(entity, options)
	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// DeleteTicket simule la suppression d'un ticket
func (m *MockGameRepository) DeleteTicket(obj *transfert.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)
	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// CountTicket simule le comptage de tickets
func (m *MockGameRepository) CountTicket(obj *transfert.Ticket, options ...database.Option) (int, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Error(1) != nil {
		return 0, args.Error(1).(errors.ErrorInterface)
	}
	return args.Int(0), nil
}

// Tests pour la méthode HydrateDBWithTickets
func TestHydrateDBWithTickets(t *testing.T) {
	// Initialisation du MockGameRepository
	mockRepo := new(MockGameRepository)

	token1 := token.Generate(12)
	token2 := token.Generate(12)

	// Configuration du mock pour ReadTickets
	mockRepo.On("ReadTickets", mock.Anything, mock.Anything).Return([]*entities.Ticket{
		{Token: token1},
		{Token: token2},
	}, errors.ErrorInterface(nil))

	// Configuration du mock pour CountTicket
	mockRepo.On("CountTicket", mock.MatchedBy(func(ticket *transfert.Ticket) bool {
		return ticket != nil && ticket.Prize != nil && *ticket.Prize == "PrizeA"
	}), mock.Anything).Return(100, errors.ErrorInterface(nil))

	mockRepo.On("CountTicket", mock.MatchedBy(func(ticket *transfert.Ticket) bool {
		return ticket != nil && ticket.Prize != nil && *ticket.Prize == "PrizeB"
	}), mock.Anything).Return(200, errors.ErrorInterface(nil))

	// Configuration du mock pour CreateTickets
	mockRepo.On("CreateTickets", mock.Anything, mock.Anything).Return(errors.ErrorInterface(nil))

	// Appel de la méthode HydrateDBWithTickets
	dispatch := map[string]int{
		"PrizeA": 50,
		"PrizeB": 50,
	}
	events.HydrateDBWithTickets(mockRepo, 1000, dispatch)

	// Vérifications
	mockRepo.AssertCalled(t, "ReadTickets", mock.Anything, mock.Anything)
	mockRepo.AssertCalled(t, "CountTicket", mock.MatchedBy(func(ticket *transfert.Ticket) bool {
		return ticket != nil && ticket.Prize != nil && *ticket.Prize == "PrizeA"
	}), mock.Anything)
	mockRepo.AssertCalled(t, "CountTicket", mock.MatchedBy(func(ticket *transfert.Ticket) bool {
		return ticket != nil && ticket.Prize != nil && *ticket.Prize == "PrizeB"
	}), mock.Anything)
	mockRepo.AssertCalled(t, "CreateTickets", mock.Anything, mock.Anything)
}
