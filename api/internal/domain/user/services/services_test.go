package services_test

import (
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) ReadUser(user *transfert.User) (*entities.Client, *entities.Employee, error) {
	args := m.Called(user)
	if args.Get(0) == nil && args.Get(1) != nil && args.Get(2) == nil {
		return nil, args.Get(1).(*entities.Employee), nil
	}

	if args.Get(0) != nil && args.Get(1) == nil && args.Get(2) == nil {
		return args.Get(0).(*entities.Client), nil, nil
	}
	return nil, nil, args.Error(2)
}

func (m *UserRepositoryMock) CreateClient(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *UserRepositoryMock) ReadClient(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *UserRepositoryMock) UpdateClient(client *entities.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *UserRepositoryMock) DeleteClient(client *transfert.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *UserRepositoryMock) CreateEmployee(employee *transfert.Employee) (*entities.Employee, error) {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Employee), args.Error(1)
}

func (m *UserRepositoryMock) ReadEmployee(employee *transfert.Employee) (*entities.Employee, error) {
	args := m.Called(employee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Employee), args.Error(1)
}

func (m *UserRepositoryMock) UpdateEmployee(employee *entities.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *UserRepositoryMock) DeleteEmployee(employee *transfert.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *UserRepositoryMock) CreateValidation(validation *transfert.Validation) (*entities.Validation, error) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (m *UserRepositoryMock) ReadValidation(validation *transfert.Validation) (*entities.Validation, error) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (m *UserRepositoryMock) UpdateValidation(validation *entities.Validation) error {
	args := m.Called(validation)

	if args.Get(0) == nil {
		validation.ID = uuid.New().String()
		validation.Token = token.NewLuhn("666666").Pointer()
		return nil
	}

	return args.Error(0)
}

func (m *UserRepositoryMock) DeleteValidation(validation *transfert.Validation) error {
	args := m.Called(validation)
	return args.Error(0)
}

func (m *UserRepositoryMock) CreateCredential(credential *transfert.Credential) (*entities.Credential, error) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credential), args.Error(1)
}

func (m *UserRepositoryMock) ReadCredential(credential *transfert.Credential) (*entities.Credential, error) {
	args := m.Called(credential)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credential), args.Error(1)
}

func (m *UserRepositoryMock) UpdateCredential(credential *entities.Credential) error {
	args := m.Called(credential)
	if args.Get(0) == nil {
		credential.ID = uuid.New().String()
		return nil
	}
	return args.Error(0)
}

func (m *UserRepositoryMock) DeleteCredential(credential *transfert.Credential) error {
	args := m.Called(credential)
	return args.Error(0)
}

type MailServiceMock struct {
	mock.Mock
}

func (m *MailServiceMock) Send(mail *mail.Mail) error {
	args := m.Called(mail)
	return args.Error(0)
}

func (m *MailServiceMock) From() string {
	args := m.Called()
	return args.String(0)
}

func (m *MailServiceMock) Expeditor() string {
	args := m.Called()
	return args.String(0)
}

func setup() (*services.UserService, *UserRepositoryMock, *MailServiceMock) {
	mockRepository := new(UserRepositoryMock)
	mockMailer := new(MailServiceMock)
	service := services.User(mockRepository, mockMailer)

	return service, mockRepository, mockMailer
}
