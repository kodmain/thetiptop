package services_test

import (
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	gameTransfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	gameEntity "github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) ReadUser(user *transfert.User, options ...database.Option) (*entities.Client, *entities.Employee, errors.ErrorInterface) {
	args := m.Called(user)
	if args.Get(0) == nil && args.Get(1) != nil && args.Get(2) == nil {
		return nil, args.Get(1).(*entities.Employee), nil
	}

	if args.Get(0) != nil && args.Get(1) == nil && args.Get(2) == nil {
		return args.Get(0).(*entities.Client), nil, nil
	}
	return nil, nil, args.Get(2).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateClient(client *transfert.Client, options ...database.Option) (*entities.Client, errors.ErrorInterface) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Client), nil
}

func (m *UserRepositoryMock) ReadClient(client *transfert.Client, options ...database.Option) (*entities.Client, errors.ErrorInterface) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Client), nil
}

func (m *UserRepositoryMock) UpdateClient(client *entities.Client, options ...database.Option) errors.ErrorInterface {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteClient(client *transfert.Client, options ...database.Option) errors.ErrorInterface {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateEmployee(employee *transfert.Employee, options ...database.Option) (*entities.Employee, errors.ErrorInterface) {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}

func (m *UserRepositoryMock) ReadEmployee(employee *transfert.Employee, options ...database.Option) (*entities.Employee, errors.ErrorInterface) {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}

func (m *UserRepositoryMock) UpdateEmployee(employee *entities.Employee, options ...database.Option) errors.ErrorInterface {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteEmployee(employee *transfert.Employee, options ...database.Option) errors.ErrorInterface {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateValidation(validation *transfert.Validation, options ...database.Option) (*entities.Validation, errors.ErrorInterface) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Validation), nil
}

func (m *UserRepositoryMock) ReadValidation(validation *transfert.Validation, options ...database.Option) (*entities.Validation, errors.ErrorInterface) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Validation), nil
}

func (m *UserRepositoryMock) ReadValidations(validation *transfert.Validation, options ...database.Option) ([]*entities.Validation, errors.ErrorInterface) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).([]*entities.Validation), nil
}

func (m *UserRepositoryMock) UpdateValidation(validation *entities.Validation, options ...database.Option) errors.ErrorInterface {
	args := m.Called(validation)

	if args.Get(0) == nil {
		validation.ID = uuid.New().String()
		validation.Token = token.NewLuhn("666666").Pointer()
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteValidation(validation *transfert.Validation, options ...database.Option) errors.ErrorInterface {
	args := m.Called(validation)
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateCredential(credential *transfert.Credential, options ...database.Option) (*entities.Credential, errors.ErrorInterface) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Credential), nil
}

func (m *UserRepositoryMock) ReadCredential(credential *transfert.Credential, options ...database.Option) (*entities.Credential, errors.ErrorInterface) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Credential), nil
}

func (m *UserRepositoryMock) UpdateCredential(credential *entities.Credential, options ...database.Option) errors.ErrorInterface {
	args := m.Called(credential)
	if args.Get(0) == nil {
		credential.ID = uuid.New().String()
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteCredential(credential *transfert.Credential, options ...database.Option) errors.ErrorInterface {
	args := m.Called(credential)
	return args.Get(0).(errors.ErrorInterface)
}

type MailServiceMock struct {
	mock.Mock
}

func (m *MailServiceMock) Send(mail *mail.Mail) error {
	args := m.Called(mail)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *MailServiceMock) From() string {
	args := m.Called()
	return args.String(0)
}

func (m *MailServiceMock) Expeditor() string {
	args := m.Called()
	return args.String(0)
}

type PermissionMock struct {
	mock.Mock
}

func (m *PermissionMock) IsAuthenticated() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *PermissionMock) GetCredentialID() *string {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*string)
}

func (m *PermissionMock) IsGrantedByRules(rules ...security.Rule) bool {
	args := m.Called(rules)
	return args.Bool(0)
}

func (m *PermissionMock) IsGrantedByRoles(roles ...security.Role) bool {
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

type GameRepositoryMock struct {
	mock.Mock
}

// CreateTicket simule la création d'un ticket.
func (m *GameRepositoryMock) CreateTicket(obj *gameTransfert.Ticket, options ...database.Option) (*gameEntity.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}

	return args.Get(0).(*gameEntity.Ticket), nil
}

// CreateTickets simule la création de plusieurs tickets.
func (m *GameRepositoryMock) CreateTickets(objs []*gameTransfert.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0).(errors.ErrorInterface)
}

// ReadTicket simule la lecture d'un ticket.
func (m *GameRepositoryMock) ReadTicket(obj *gameTransfert.Ticket, options ...database.Option) (*gameEntity.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}

	return args.Get(0).(*gameEntity.Ticket), nil
}

// ReadTickets simule la lecture de plusieurs tickets.
func (m *GameRepositoryMock) ReadTickets(obj *gameTransfert.Ticket, options ...database.Option) ([]*gameEntity.Ticket, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}

	return args.Get(0).([]*gameEntity.Ticket), nil
}

// UpdateTicket simule la mise à jour d'un ticket.
func (m *GameRepositoryMock) UpdateTicket(entity *gameEntity.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(entity, options)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0).(errors.ErrorInterface)
}

// DeleteTicket simule la suppression d'un ticket.
func (m *GameRepositoryMock) DeleteTicket(obj *gameTransfert.Ticket, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil
	}

	return args.Error(0).(errors.ErrorInterface)
}

// CountTicket simule le comptage des tickets.
func (m *GameRepositoryMock) CountTicket(obj *gameTransfert.Ticket, options ...database.Option) (int, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return 0, args.Error(1).(errors.ErrorInterface)
	}

	return args.Int(0), nil
}

func setup() (*services.UserService, *UserRepositoryMock, *MailServiceMock, *PermissionMock, *GameRepositoryMock) {
	mockRepository := new(UserRepositoryMock)
	gameRepository := new(GameRepositoryMock)
	mockMailer := new(MailServiceMock)
	mockSecurity := new(PermissionMock)
	service := services.User(mockSecurity, mockRepository, gameRepository, mockMailer)

	return service, mockRepository, mockMailer, mockSecurity, gameRepository
}
