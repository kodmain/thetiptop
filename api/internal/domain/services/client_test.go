package services_test

import (
	"fmt"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ClientRepositoryMock struct {
	mock.Mock
}

func (m *ClientRepositoryMock) Create(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *ClientRepositoryMock) Read(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
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

func setup() (*services.ClientService, *ClientRepositoryMock, *MailServiceMock) {
	mockRepository := new(ClientRepositoryMock)
	mockMailer := new(MailServiceMock)
	service := services.Client(mockRepository, mockMailer)

	return service, mockRepository, mockMailer
}

func TestSignUp(t *testing.T) {

	t.Run("nil input", func(t *testing.T) {
		service, _, _ := setup()
		require.NotNil(t, service)

		result, err := service.SignUp(nil)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrNoDto, err.Error())
	})

	inputClient := &transfert.Client{
		Email:    "hello@thetiptop",
		Password: "azertyuiop",
	}

	expectedClient := &entities.Client{
		ID:       "42debee6-2063-4566-baf1-37a7bdd139ff",
		Email:    "hello@thetiptop",
		Password: "$2a$10$wO5PfDAGp6w2ubKp0vEdXeUe2HlfOv5iRJ3C3MVR0vJhscD0G.NKS", // hashed password
	}

	t.Run("client already exists", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		result, err := service.SignUp(inputClient)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientAlreadyExists, err.Error())
		mockRepository.AssertExpectations(t)
	})

	t.Run("failed signup", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()

		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockRepository.On("Create", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf("failed to create client"))

		result, err := service.SignUp(inputClient)
		require.Error(t, err)
		require.Nil(t, result)

		mockRepository.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("successful signup", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()

		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockRepository.On("Create", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)

		result, err := service.SignUp(inputClient)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expectedClient.Email, result.Email)

		mockRepository.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("failed mail send", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()

		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockRepository.On("Create", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(fmt.Errorf("failed to send mail"))

		result, err := service.SignUp(inputClient)
		require.Error(t, err)
		require.Nil(t, result)

		mockRepository.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	inputClient := &transfert.Client{
		Email:    "hello@thetiptop",
		Password: "azertyuiop",
	}

	hashedPassword, err := hash.Hash(inputClient.Email+":"+inputClient.Password, hash.BCRYPT)
	require.NoError(t, err)

	expectedClient := &entities.Client{
		ID:       "42debee6-2063-4566-baf1-37a7bdd139ff",
		Email:    inputClient.Email,
		Password: hashedPassword,
	}

	t.Run("client not found", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		result, err := service.SignIn(inputClient)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientNotFound, err.Error())
		mockRepository.AssertExpectations(t)
	})

	t.Run("incorrect password", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)

		wrongPasswordClient := &transfert.Client{
			Email:    "hello@thetiptop",
			Password: "wrongpassword",
		}

		result, err := service.SignIn(wrongPasswordClient)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientNotFound, err.Error())
		mockRepository.AssertExpectations(t)
	})

	t.Run("successful signin", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("Read", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)

		result, err := service.SignIn(inputClient)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expectedClient.Email, result.Email)
		mockRepository.AssertExpectations(t)
	})
}
