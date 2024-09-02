package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ClientRepositoryMock struct {
	mock.Mock
}

func (m *ClientRepositoryMock) CreateClient(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *ClientRepositoryMock) ReadClient(client *transfert.Client) (*entities.Client, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (m *ClientRepositoryMock) UpdateClient(client *entities.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientRepositoryMock) DeleteClient(client *transfert.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientRepositoryMock) CreateValidation(validation *transfert.Validation) (*entities.Validation, error) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (m *ClientRepositoryMock) ReadValidation(validation *transfert.Validation) (*entities.Validation, error) {
	args := m.Called(validation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (m *ClientRepositoryMock) UpdateValidation(validation *entities.Validation) error {
	args := m.Called(validation)

	if args.Get(0) == nil {
		validation.ID = uuid.New().String()
		validation.Token = token.NewLuhn("666666").Pointer()
		return nil
	}

	return args.Error(0)
}

func (m *ClientRepositoryMock) DeleteValidation(validation *transfert.Validation) error {
	args := m.Called(validation)
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

	idClient, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)
	idValidation, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	inputClient := &transfert.Client{
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("azertyuiop"),
	}

	sidClient := idClient.String()

	expectedClient := &entities.Client{
		ID:       idClient.String(),
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("$2a$10$wO5PfDAGp6w2ubKp0vEdXeUe2HlfOv5iRJ3C3MVR0vJhscD0G.NKS"), // hashed password
		Validations: []*entities.Validation{
			{
				ID:        idValidation.String(),
				Token:     token.NewLuhn("666666").Pointer(),
				Type:      0,
				Validated: false,
				ClientID:  &sidClient,
			},
		},
	}

	t.Run("client already exists", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		result, err := service.SignUp(inputClient)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientAlreadyExists, err.Error())
		mockRepository.AssertExpectations(t)
	})

	t.Run("failed signup", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()

		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockRepository.On("CreateClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf("failed to create client"))

		result, err := service.SignUp(inputClient)
		require.Error(t, err)
		require.Nil(t, result)

		mockRepository.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("successful signup", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()

		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockRepository.On("CreateClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(nil)
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)

		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)

		result, err := service.SignUp(inputClient)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expectedClient.Email, result.Email)

		time.Sleep(1 * time.Second)

		mockRepository.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})

	t.Run("failed mail send", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()

		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockRepository.On("CreateClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(nil)
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(fmt.Errorf("failed to send mail"))

		result, err := service.SignUp(inputClient)
		require.NoError(t, err)
		require.NotNil(t, result)

		time.Sleep(100 * time.Millisecond)

		mockRepository.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	idClient, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)
	idValidation, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	inputClient := &transfert.Client{
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("azertyuiop"),
	}

	hashedPassword, err := hash.Hash(aws.String(*inputClient.Email+":"+*inputClient.Password), hash.BCRYPT)
	require.NoError(t, err)

	sidClient := idClient.String()

	expectedClient := &entities.Client{
		ID:       idClient.String(),
		Email:    inputClient.Email,
		Password: hashedPassword,
		Validations: []*entities.Validation{
			{
				ID:        idValidation.String(),
				Token:     token.NewLuhn("666666").Pointer(),
				Type:      0,
				Validated: true,
				ClientID:  &sidClient,
			},
		},
	}

	t.Run("client not found", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		result, err := service.SignIn(inputClient)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientNotFound, err.Error())
		mockRepository.AssertExpectations(t)
	})

	t.Run("incorrect password", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)

		wrongPasswordClient := &transfert.Client{
			Email:    aws.String("hello@thetiptop"),
			Password: aws.String("wrongpassword"),
		}

		result, err := service.SignIn(wrongPasswordClient)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrClientNotFound, err.Error())
		mockRepository.AssertExpectations(t)
	})

	t.Run("successful signin", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)

		result, err := service.SignIn(inputClient)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expectedClient.Email, result.Email)
		mockRepository.AssertExpectations(t)
	})
}

func TestValidationRecover(t *testing.T) {

	idClient, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	inputClient := &transfert.Client{
		Email: aws.String("hello@thetiptop"),
	}

	expectedClient := &entities.Client{
		ID:    idClient.String(),
		Email: inputClient.Email,
		Validations: []*entities.Validation{
			{
				ID:        uuid.New().String(),
				Type:      entities.MailValidation,
				Validated: true,
				ClientID:  aws.String(idClient.String()),
			},
		},
	}

	t.Run("validation recover", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(nil)
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)
		err := service.ValidationRecover(inputClient)
		require.NoError(t, err)
	})

	t.Run("client not found", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)
		err := service.ValidationRecover(inputClient)
		require.Error(t, err)
	})

	t.Run("client not validated", func(t *testing.T) {
		service, mockRepository, _ := setup()
		clientWithoutValidation := &entities.Client{
			ID:          idClient.String(),
			Email:       inputClient.Email,
			Validations: []*entities.Validation{},
		}
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(clientWithoutValidation, nil)
		err := service.ValidationRecover(inputClient)
		require.Error(t, err)
		require.Equal(t, fmt.Errorf(errors.ErrClientNotValidate, entities.MailValidation.String()), err)
	})

	t.Run("validation update fail", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(fmt.Errorf("failed to update validation"))
		err := service.ValidationRecover(inputClient)
		require.Error(t, err)
	})

	t.Run("client update fail", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(fmt.Errorf("failed to update client"))
		err := service.ValidationRecover(inputClient)
		require.Error(t, err)
	})
}

func TestPasswordRecover(t *testing.T) {

	idClient, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)
	idValidation, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	inputClient := &transfert.Client{
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("azertyuiop"),
	}

	hashedPassword, err := hash.Hash(aws.String(*inputClient.Email+":"+*inputClient.Password), hash.BCRYPT)
	require.NoError(t, err)

	sidClient := idClient.String()

	expectedClient := &entities.Client{
		ID:       idClient.String(),
		Email:    inputClient.Email,
		Password: hashedPassword,
		Validations: []*entities.Validation{
			{
				ID:        idValidation.String(),
				Token:     token.NewLuhn("666666").Pointer(),
				Type:      0,
				Validated: true,
				ClientID:  &sidClient,
			},
		},
	}
	t.Run("password recover", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(nil)
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)
		err := service.PasswordRecover(inputClient)
		require.NoError(t, err)
	})

	t.Run("client not found", func(t *testing.T) {
		service, mockRepository, mockMailer := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)
		err := service.PasswordRecover(inputClient)
		require.Error(t, err)
	})

	t.Run("user update fail", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(fmt.Errorf("failed to update client"))
		mockRepository.On("UpdateValidation", mock.AnythingOfType("*entities.Validation")).Return(nil)
		err := service.PasswordRecover(inputClient)
		require.Error(t, err)
	})

	expectedClient.Validations[0].Validated = false

	t.Run("user fail validation", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		err := service.PasswordRecover(inputClient)
		require.Error(t, err)
	})

}

func TestPasswordUpdate(t *testing.T) {

	idClient, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)
	idValidation, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	sidClient := idClient.String()

	inputClient := &transfert.Client{
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("azertyuiop"),
	}

	hashedPassword, err := hash.Hash(aws.String(*inputClient.Email+":"+*inputClient.Password), hash.BCRYPT)
	require.NoError(t, err)

	expectedClient := &entities.Client{
		ID:       sidClient,
		Email:    inputClient.Email,
		Password: hashedPassword,
		Validations: []*entities.Validation{
			{
				ID:        idValidation.String(),
				Token:     token.NewLuhn("555555").Pointer(),
				Type:      entities.MailValidation,
				Validated: true,
				ClientID:  &sidClient,
			},
			{
				ID:        idValidation.String(),
				Token:     token.NewLuhn("666666").Pointer(),
				Type:      entities.PasswordRecover,
				Validated: true,
				ClientID:  &sidClient,
			},
		},
	}

	t.Run("successful update user", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(nil)
		err := service.PasswordUpdate(inputClient)
		require.NoError(t, err)
	})

	t.Run("user not exist", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		err := service.PasswordUpdate(inputClient)
		require.Error(t, err)
	})

	t.Run("user not exist", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
		err := service.PasswordUpdate(inputClient)
		require.Error(t, err)
	})

	t.Run("user update fail", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		mockRepository.On("UpdateClient", mock.AnythingOfType("*entities.Client")).Return(fmt.Errorf("failed to update client"))
		err := service.PasswordUpdate(inputClient)
		require.Error(t, err)
	})

	expectedClient.Validations[0].Validated = false

	t.Run("user fail password validation", func(t *testing.T) {
		service, mockRepository, _ := setup()
		mockRepository.On("ReadClient", mock.AnythingOfType("*transfert.Client")).Return(expectedClient, nil)
		err := service.PasswordUpdate(inputClient)
		require.Error(t, err)
	})
}
