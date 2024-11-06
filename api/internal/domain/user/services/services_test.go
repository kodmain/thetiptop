package services_test

import (
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) ReadUser(user *transfert.User) (*entities.Client, *entities.Employee, errors.ErrorInterface) {
	args := m.Called(user)
	if args.Get(0) == nil && args.Get(1) != nil && args.Get(2) == nil {
		return nil, args.Get(1).(*entities.Employee), nil
	}

	if args.Get(0) != nil && args.Get(1) == nil && args.Get(2) == nil {
		return args.Get(0).(*entities.Client), nil, nil
	}
	return nil, nil, args.Get(2).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateClient(client *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Client), nil
}

func (m *UserRepositoryMock) ReadClient(client *transfert.Client) (*entities.Client, errors.ErrorInterface) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Client), nil
}

func (m *UserRepositoryMock) UpdateClient(client *entities.Client) errors.ErrorInterface {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteClient(client *transfert.Client) errors.ErrorInterface {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateEmployee(employee *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}

func (m *UserRepositoryMock) ReadEmployee(employee *transfert.Employee) (*entities.Employee, errors.ErrorInterface) {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Employee), nil
}

func (m *UserRepositoryMock) UpdateEmployee(employee *entities.Employee) errors.ErrorInterface {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteEmployee(employee *transfert.Employee) errors.ErrorInterface {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateValidation(validation *transfert.Validation) (*entities.Validation, errors.ErrorInterface) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Validation), nil
}

func (m *UserRepositoryMock) ReadValidation(validation *transfert.Validation) (*entities.Validation, errors.ErrorInterface) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Validation), nil
}

func (m *UserRepositoryMock) UpdateValidation(validation *entities.Validation) errors.ErrorInterface {
	args := m.Called(validation)

	if args.Get(0) == nil {
		validation.ID = uuid.New().String()
		validation.Token = token.NewLuhn("666666").Pointer()
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteValidation(validation *transfert.Validation) errors.ErrorInterface {
	args := m.Called(validation)
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) CreateCredential(credential *transfert.Credential) (*entities.Credential, errors.ErrorInterface) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Credential), nil
}

func (m *UserRepositoryMock) ReadCredential(credential *transfert.Credential) (*entities.Credential, errors.ErrorInterface) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Get(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Credential), nil
}

func (m *UserRepositoryMock) UpdateCredential(credential *entities.Credential) errors.ErrorInterface {
	args := m.Called(credential)
	if args.Get(0) == nil {
		credential.ID = uuid.New().String()
		return nil
	}
	return args.Get(0).(errors.ErrorInterface)
}

func (m *UserRepositoryMock) DeleteCredential(credential *transfert.Credential) errors.ErrorInterface {
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

func (m *PermissionMock) IsGranted(roles ...string) bool {
	args := m.Called(roles)
	return args.Bool(0)
}

func (m *PermissionMock) CanRead(ressource entities.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanCreate(ressource entities.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanUpdate(ressource entities.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanDelete(ressource entities.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func setup() (*services.UserService, *UserRepositoryMock, *MailServiceMock, *PermissionMock) {
	mockRepository := new(UserRepositoryMock)
	mockMailer := new(MailServiceMock)
	mockSecurity := new(PermissionMock)
	service := services.User(mockSecurity, mockRepository, mockMailer)

	return service, mockRepository, mockMailer, mockSecurity
}
