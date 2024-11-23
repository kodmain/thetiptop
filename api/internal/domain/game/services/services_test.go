package services_test

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/game/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/mock"
)

// GameRepositoryMock est le mock pour GameRepositoryInterface
type GameRepositoryMock struct {
	mock.Mock
}

// CreateTicket simule la création d'un ticket.
func (m *GameRepositoryMock) CreateTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}

	return args.Get(0).(*entities.Ticket), nil
}

// CreateTickets simule la création de plusieurs tickets.
func (m *GameRepositoryMock) CreateTickets(objs []*transfert.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0).(errors.ErrorInterface)
}

// ReadTicket simule la lecture d'un ticket.
func (m *GameRepositoryMock) ReadTicket(obj *transfert.Ticket, options ...database.Option) (*entities.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}

	return args.Get(0).(*entities.Ticket), nil
}

// ReadTickets simule la lecture de plusieurs tickets.
func (m *GameRepositoryMock) ReadTickets(obj *transfert.Ticket, options ...database.Option) ([]*entities.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}

	return args.Get(0).([]*entities.Ticket), nil
}

// UpdateTicket simule la mise à jour d'un ticket.
func (m *GameRepositoryMock) UpdateTicket(entity *entities.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(entity, options)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0).(errors.ErrorInterface)
}

// DeleteTicket simule la suppression d'un ticket.
func (m *GameRepositoryMock) DeleteTicket(obj *transfert.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0).(errors.ErrorInterface)
}

// CountTicket simule le comptage des tickets.
func (m *GameRepositoryMock) CountTicket(obj *transfert.Ticket, options ...database.Option) (int, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return 0, args.Error(1).(errors.ErrorInterface)
	}

	return args.Int(0), nil
}

// PermissionMock est le mock pour PermissionInterface
type PermissionMock struct {
	mock.Mock
}

func (m *PermissionMock) IsAuthenticated() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *PermissionMock) IsGrantedByRoles(roles ...security.Role) bool {
	args := m.Called(roles)
	return args.Bool(0)
}

func (m *PermissionMock) IsGrantedByRules(roles ...security.Rule) bool {
	args := m.Called(roles)
	return args.Bool(0)
}

func (m *PermissionMock) CanRead(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanCreate(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanUpdate(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanDelete(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) GetCredentialID() *string {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*string)
}

func setup() (*services.GameService, *GameRepositoryMock, *PermissionMock) {
	mockRepository := new(GameRepositoryMock)
	mockSecurity := new(PermissionMock)

	service := services.Game(mockSecurity, mockRepository)

	return service, mockRepository, mockSecurity
}
